-- Migration: Add Live Chat Support to Messages Table
-- Date: 2024-12-04
-- Description: Extends messages table to support both private 1-on-1 chat and TikTok-style live streaming chat

-- ==================== STEP 1: Add Live Chat Columns ====================

ALTER TABLE messages 
ADD COLUMN IF NOT EXISTS live_stream_id UUID REFERENCES users(id) ON DELETE CASCADE,
ADD COLUMN IF NOT EXISTS is_live BOOLEAN DEFAULT FALSE NOT NULL,
ADD COLUMN IF NOT EXISTS is_system BOOLEAN DEFAULT FALSE NOT NULL,
ADD COLUMN IF NOT EXISTS seq BIGINT DEFAULT 0,
ADD COLUMN IF NOT EXISTS pinned BOOLEAN DEFAULT FALSE NOT NULL;

-- ==================== STEP 2: Add Indexes for Performance ====================

-- Index for live chat queries (fetch messages by stream)
CREATE INDEX IF NOT EXISTS idx_messages_live_stream_id 
ON messages(live_stream_id) 
WHERE is_live = TRUE;

-- Index for sequence-based queries (reconnection & replay)
CREATE INDEX IF NOT EXISTS idx_messages_live_stream_seq 
ON messages(live_stream_id, seq) 
WHERE is_live = TRUE;

-- Index for pinned messages
CREATE INDEX IF NOT EXISTS idx_messages_pinned 
ON messages(live_stream_id, pinned) 
WHERE is_live = TRUE AND pinned = TRUE;

-- Index for system messages
CREATE INDEX IF NOT EXISTS idx_messages_system 
ON messages(is_system, live_stream_id) 
WHERE is_system = TRUE;

-- Composite index for live chat pagination
CREATE INDEX IF NOT EXISTS idx_messages_live_created 
ON messages(live_stream_id, created_at DESC) 
WHERE is_live = TRUE;

-- ==================== STEP 3: Update Constraints ====================

-- Make match_id nullable (since live messages don't have matches)
ALTER TABLE messages 
ALTER COLUMN match_id DROP NOT NULL;

-- Make receiver_id nullable (since live messages are broadcast)
ALTER TABLE messages 
ALTER COLUMN receiver_id DROP NOT NULL;

-- Add check constraint: message must be either private OR live
ALTER TABLE messages 
ADD CONSTRAINT chk_message_mode 
CHECK (
    (match_id IS NOT NULL AND receiver_id IS NOT NULL AND is_live = FALSE AND live_stream_id IS NULL)
    OR
    (live_stream_id IS NOT NULL AND is_live = TRUE AND match_id IS NULL)
);

-- ==================== STEP 4: Add Comments ====================

COMMENT ON COLUMN messages.live_stream_id IS 'Reference to live stream (user who is broadcasting). NULL for private messages.';
COMMENT ON COLUMN messages.is_live IS 'TRUE for live streaming messages, FALSE for private 1-on-1 messages';
COMMENT ON COLUMN messages.is_system IS 'TRUE for system messages (join/leave/announcements). Only broadcaster can send.';
COMMENT ON COLUMN messages.seq IS 'Monotonic sequence number for live messages. Used for ordering and replay. 0 for private messages.';
COMMENT ON COLUMN messages.pinned IS 'TRUE if message is pinned by broadcaster. Only one pinned message per stream.';

-- ==================== STEP 5: Create Live Streams Table (Optional) ====================

CREATE TABLE IF NOT EXISTS live_streams (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL DEFAULT 'Live Stream',
    description TEXT,
    thumbnail_url TEXT,
    
    -- Streaming details
    stream_key VARCHAR(255) UNIQUE NOT NULL,
    rtmp_url TEXT,
    playback_url TEXT,
    
    -- Status
    status VARCHAR(50) DEFAULT 'pending' CHECK (status IN ('pending', 'live', 'ended', 'banned')),
    started_at TIMESTAMPTZ,
    ended_at TIMESTAMPTZ,
    
    -- Stats
    peak_viewers INT DEFAULT 0,
    total_views INT DEFAULT 0,
    total_messages INT DEFAULT 0,
    total_gifts_received INT DEFAULT 0,
    total_coins_earned INT DEFAULT 0,
    
    -- Settings
    allow_chat BOOLEAN DEFAULT TRUE,
    allow_gifts BOOLEAN DEFAULT TRUE,
    is_private BOOLEAN DEFAULT FALSE,
    
    -- Metadata
    metadata JSONB DEFAULT '{}',
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes for live_streams
CREATE INDEX IF NOT EXISTS idx_live_streams_user_id ON live_streams(user_id);
CREATE INDEX IF NOT EXISTS idx_live_streams_status ON live_streams(status);
CREATE INDEX IF NOT EXISTS idx_live_streams_started_at ON live_streams(started_at DESC) WHERE status = 'live';
CREATE UNIQUE INDEX IF NOT EXISTS idx_live_streams_stream_key ON live_streams(stream_key);

-- ==================== STEP 6: Create Viewer Tracking Table (Optional) ====================

CREATE TABLE IF NOT EXISTS live_stream_viewers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    live_stream_id UUID NOT NULL REFERENCES live_streams(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    joined_at TIMESTAMPTZ DEFAULT NOW(),
    left_at TIMESTAMPTZ,
    
    -- Stats
    total_messages_sent INT DEFAULT 0,
    total_gifts_sent INT DEFAULT 0,
    total_watch_time_seconds INT DEFAULT 0,
    
    UNIQUE(live_stream_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_live_stream_viewers_stream ON live_stream_viewers(live_stream_id);
CREATE INDEX IF NOT EXISTS idx_live_stream_viewers_user ON live_stream_viewers(user_id);
CREATE INDEX IF NOT EXISTS idx_live_stream_viewers_active ON live_stream_viewers(live_stream_id, left_at) WHERE left_at IS NULL;

-- ==================== STEP 7: Update Existing Data ====================

-- Set is_live = FALSE for all existing messages (they are private messages)
UPDATE messages 
SET is_live = FALSE, 
    is_system = FALSE, 
    seq = 0, 
    pinned = FALSE
WHERE is_live IS NULL;

-- ==================== STEP 8: Create Helper Functions ====================

-- Function to get current viewer count for a live stream
CREATE OR REPLACE FUNCTION get_live_viewer_count(stream_id UUID)
RETURNS INT AS $$
    SELECT COUNT(*)::INT
    FROM live_stream_viewers
    WHERE live_stream_id = stream_id 
    AND left_at IS NULL;
$$ LANGUAGE SQL STABLE;

-- Function to update live stream stats
CREATE OR REPLACE FUNCTION update_live_stream_stats()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_live = TRUE AND NEW.live_stream_id IS NOT NULL THEN
        UPDATE live_streams
        SET total_messages = total_messages + 1,
            updated_at = NOW()
        WHERE id = NEW.live_stream_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to auto-update live stream stats on new message
CREATE TRIGGER trg_update_live_stream_stats
AFTER INSERT ON messages
FOR EACH ROW
WHEN (NEW.is_live = TRUE)
EXECUTE FUNCTION update_live_stream_stats();

-- ==================== STEP 9: Grant Permissions ====================

-- Grant permissions (adjust role name as needed)
-- GRANT SELECT, INSERT, UPDATE ON messages TO your_app_role;
-- GRANT SELECT, INSERT, UPDATE, DELETE ON live_streams TO your_app_role;
-- GRANT SELECT, INSERT, UPDATE, DELETE ON live_stream_viewers TO your_app_role;

-- ==================== ROLLBACK SCRIPT (for reference) ====================

/*
-- To rollback this migration:

DROP TRIGGER IF EXISTS trg_update_live_stream_stats ON messages;
DROP FUNCTION IF EXISTS update_live_stream_stats();
DROP FUNCTION IF EXISTS get_live_viewer_count(UUID);

DROP TABLE IF EXISTS live_stream_viewers CASCADE;
DROP TABLE IF EXISTS live_streams CASCADE;

ALTER TABLE messages DROP CONSTRAINT IF EXISTS chk_message_mode;
ALTER TABLE messages ALTER COLUMN match_id SET NOT NULL;
ALTER TABLE messages ALTER COLUMN receiver_id SET NOT NULL;

DROP INDEX IF EXISTS idx_messages_live_created;
DROP INDEX IF EXISTS idx_messages_system;
DROP INDEX IF EXISTS idx_messages_pinned;
DROP INDEX IF EXISTS idx_messages_live_stream_seq;
DROP INDEX IF EXISTS idx_messages_live_stream_id;

ALTER TABLE messages 
DROP COLUMN IF EXISTS pinned,
DROP COLUMN IF EXISTS seq,
DROP COLUMN IF EXISTS is_system,
DROP COLUMN IF EXISTS is_live,
DROP COLUMN IF EXISTS live_stream_id;
*/
