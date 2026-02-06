ALTER TABLE reviews 
ADD COLUMN moderation_status TEXT NOT NULL DEFAULT 'visible';

ALTER TABLE reviews 
ADD COLUMN report_count INTEGER NOT NULL DEFAULT 0;