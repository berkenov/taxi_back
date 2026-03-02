-- OTP codes for auth (temporary, expire after 5 minutes)
CREATE TABLE IF NOT EXISTS otp_codes (
    phone VARCHAR(20) NOT NULL PRIMARY KEY,
    code VARCHAR(4) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_otp_codes_expires_at ON otp_codes(expires_at);
