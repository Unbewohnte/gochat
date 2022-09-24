package api

// How many characters can message hold
const MaxMessageContentLen uint = 500

// How many messages can be "remembered" until removal
const MaxMessagesRemembered uint = 50

// Max filesize to accept as attachment
const MaxAttachmentSize uint64 = 104857600 // 100MB
