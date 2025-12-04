-- ============================================================================
-- DRAFT VIDEOS
-- Allow users to save videos as drafts before publishing
-- ============================================================================

-- Add draft status to videos table
ALTER TABLE videos ADD COLUMN IF NOT EXISTS is_draft BOOLEAN DEFAULT FALSE;
ALTER TABLE videos ADD COLUMN IF NOT EXISTS draft_saved_at TIMESTAMP WITH TIME ZONE;

CREATE INDEX IF NOT EXISTS idx_videos_is_draft ON videos(is_draft) WHERE is_draft = TRUE;
CREATE INDEX IF NOT EXISTS idx_videos_user_draft ON videos(user_id, is_draft) WHERE is_draft = TRUE;

COMMENT ON COLUMN videos.is_draft IS 'Whether this video is a draft (not published)';
COMMENT ON COLUMN videos.draft_saved_at IS 'When the draft was last saved';
