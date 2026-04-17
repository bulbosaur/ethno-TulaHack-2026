CREATE TABLE IF NOT EXISTS quests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug VARCHAR(100) UNIQUE NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    cover_url VARCHAR(500),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS quest_steps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    quest_id UUID NOT NULL REFERENCES quests(id) ON DELETE CASCADE,
    step_id VARCHAR(50) NOT NULL,
    step_order INTEGER NOT NULL,
    step_type VARCHAR(30) NOT NULL,
    title VARCHAR(255) NOT NULL,
    content JSONB NOT NULL,
    on_success JSONB,
    UNIQUE(quest_id, step_id)
);

CREATE TABLE IF NOT EXISTS user_quest_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    quest_id UUID NOT NULL REFERENCES quests(id) ON DELETE CASCADE,
    current_step_id VARCHAR(50),
    completed_steps JSONB DEFAULT '[]'::jsonb,
    status VARCHAR(20) DEFAULT 'in_progress',
    started_at TIMESTAMPTZ DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    UNIQUE(user_id, quest_id)
);

CREATE TABLE IF NOT EXISTS user_rewards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    quest_id UUID REFERENCES quests(id) ON DELETE SET NULL,
    reward_type VARCHAR(30) NOT NULL,
    reward_key VARCHAR(100) NOT NULL,
    granted_at TIMESTAMPTZ DEFAULT NOW(),
    metadata JSONB
);

CREATE INDEX idx_quests_slug ON quests(slug);
CREATE INDEX idx_quest_steps_quest ON quest_steps(quest_id, step_order);
CREATE INDEX idx_user_progress_user ON user_quest_progress(user_id, status);