DO $$ BEGIN
    CREATE TYPE user_role AS ENUM ('admin', 'agent');
    CREATE TYPE conversation_status AS ENUM ('open', 'assigned', 'closed');
    CREATE TYPE message_sender_type AS ENUM ('customer', 'agent');
    CREATE TYPE ticket_status AS ENUM ('open', 'in_progress', 'resolved', 'closed');
    CREATE TYPE ticket_priority AS ENUM ('low', 'medium', 'high');
EXCEPTION WHEN duplicate_object THEN null;
END $$;