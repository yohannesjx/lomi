-- ============================================================================
-- SOCIAL CONTENT SYSTEM (TikTok-style)
-- Videos, Likes, Reposts, Favorites, Comments
-- ============================================================================

-- ============================================================================
-- VIDEOS TABLE
-- ============================================================================

CREATE TABLE IF NOT EXISTS videos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Video content
    video_url TEXT NOT NULL,
    thumbnail_url TEXT,
    duration_seconds INTEGER NOT NULL,
    
    -- Metadata
    title TEXT,
    description TEXT,
    hashtags TEXT[], -- Array of hashtags
    
    -- Privacy
    is_private BOOLEAN DEFAULT FALSE,
    allow_comments BOOLEAN DEFAULT TRUE,
    allow_duet BOOLEAN DEFAULT TRUE,
    allow_stitch BOOLEAN DEFAULT TRUE,
    
    -- Engagement metrics
    views_count INTEGER DEFAULT 0,
    likes_count INTEGER DEFAULT 0,
    comments_count INTEGER DEFAULT 0,
    shares_count INTEGER DEFAULT 0,
    
    -- Moderation
    is_approved BOOLEAN DEFAULT FALSE,
    moderation_status VARCHAR(20) DEFAULT 'pending',
    moderation_notes TEXT,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    -- Indexes
    CHECK (duration_seconds > 0 AND duration_seconds <= 180) -- Max 3 minutes
);

CREATE INDEX IF NOT EXISTS idx_videos_user_id ON videos(user_id);
CREATE INDEX IF NOT EXISTS idx_videos_created_at ON videos(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_videos_is_approved ON videos(is_approved);
CREATE INDEX IF NOT EXISTS idx_videos_deleted_at ON videos(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_videos_hashtags ON videos USING GIN(hashtags);

-- ============================================================================
-- VIDEO LIKES TABLE
-- ============================================================================

CREATE TABLE IF NOT EXISTS video_likes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    video_id UUID NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, video_id)
);

CREATE INDEX IF NOT EXISTS idx_video_likes_user_id ON video_likes(user_id);
CREATE INDEX IF NOT EXISTS idx_video_likes_video_id ON video_likes(video_id);
CREATE INDEX IF NOT EXISTS idx_video_likes_created_at ON video_likes(created_at DESC);

-- ============================================================================
-- VIDEO REPOSTS TABLE
-- ============================================================================

CREATE TABLE IF NOT EXISTS video_reposts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    video_id UUID NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    caption TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, video_id)
);

CREATE INDEX IF NOT EXISTS idx_video_reposts_user_id ON video_reposts(user_id);
CREATE INDEX IF NOT EXISTS idx_video_reposts_video_id ON video_reposts(video_id);
CREATE INDEX IF NOT EXISTS idx_video_reposts_created_at ON video_reposts(created_at DESC);

-- ============================================================================
-- VIDEO FAVORITES TABLE
-- ============================================================================

CREATE TABLE IF NOT EXISTS video_favorites (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    video_id UUID NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, video_id)
);

CREATE INDEX IF NOT EXISTS idx_video_favorites_user_id ON video_favorites(user_id);
CREATE INDEX IF NOT EXISTS idx_video_favorites_video_id ON video_favorites(video_id);
CREATE INDEX IF NOT EXISTS idx_video_favorites_created_at ON video_favorites(created_at DESC);

-- ============================================================================
-- VIDEO COMMENTS TABLE
-- ============================================================================

CREATE TABLE IF NOT EXISTS video_comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    video_id UUID NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    parent_comment_id UUID REFERENCES video_comments(id) ON DELETE CASCADE,
    
    content TEXT NOT NULL,
    likes_count INTEGER DEFAULT 0,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_video_comments_video_id ON video_comments(video_id);
CREATE INDEX IF NOT EXISTS idx_video_comments_user_id ON video_comments(user_id);
CREATE INDEX IF NOT EXISTS idx_video_comments_parent_id ON video_comments(parent_comment_id);
CREATE INDEX IF NOT EXISTS idx_video_comments_created_at ON video_comments(created_at DESC);

-- ============================================================================
-- FUNCTIONS & TRIGGERS
-- ============================================================================

-- Function to update video counts
CREATE OR REPLACE FUNCTION update_video_counts()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_TABLE_NAME = 'video_likes' THEN
        IF TG_OP = 'INSERT' THEN
            UPDATE videos SET likes_count = likes_count + 1 WHERE id = NEW.video_id;
        ELSIF TG_OP = 'DELETE' THEN
            UPDATE videos SET likes_count = GREATEST(0, likes_count - 1) WHERE id = OLD.video_id;
        END IF;
    ELSIF TG_TABLE_NAME = 'video_comments' THEN
        IF TG_OP = 'INSERT' THEN
            UPDATE videos SET comments_count = comments_count + 1 WHERE id = NEW.video_id;
        ELSIF TG_OP = 'DELETE' THEN
            UPDATE videos SET comments_count = GREATEST(0, comments_count - 1) WHERE id = OLD.video_id;
        END IF;
    ELSIF TG_TABLE_NAME = 'video_reposts' THEN
        IF TG_OP = 'INSERT' THEN
            UPDATE videos SET shares_count = shares_count + 1 WHERE id = NEW.video_id;
        ELSIF TG_OP = 'DELETE' THEN
            UPDATE videos SET shares_count = GREATEST(0, shares_count - 1) WHERE id = OLD.video_id;
        END IF;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Triggers for video counts
DROP TRIGGER IF EXISTS trigger_update_video_likes_count ON video_likes;
CREATE TRIGGER trigger_update_video_likes_count
AFTER INSERT OR DELETE ON video_likes
FOR EACH ROW EXECUTE FUNCTION update_video_counts();

DROP TRIGGER IF EXISTS trigger_update_video_comments_count ON video_comments;
CREATE TRIGGER trigger_update_video_comments_count
AFTER INSERT OR DELETE ON video_comments
FOR EACH ROW EXECUTE FUNCTION update_video_counts();

DROP TRIGGER IF EXISTS trigger_update_video_shares_count ON video_reposts;
CREATE TRIGGER trigger_update_video_shares_count
AFTER INSERT OR DELETE ON video_reposts
FOR EACH ROW EXECUTE FUNCTION update_video_counts();

-- Apply updated_at trigger to videos and comments
DROP TRIGGER IF EXISTS update_videos_updated_at ON videos;
CREATE TRIGGER update_videos_updated_at 
BEFORE UPDATE ON videos 
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_video_comments_updated_at ON video_comments;
CREATE TRIGGER update_video_comments_updated_at 
BEFORE UPDATE ON video_comments 
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- COMMENTS
-- ============================================================================

COMMENT ON TABLE videos IS 'User-generated video content (TikTok-style)';
COMMENT ON TABLE video_likes IS 'Video likes from users';
COMMENT ON TABLE video_reposts IS 'Video reposts/shares';
COMMENT ON TABLE video_favorites IS 'User favorite videos';
COMMENT ON TABLE video_comments IS 'Comments on videos with threading support';
