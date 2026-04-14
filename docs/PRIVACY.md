# Privacy policy

This document describes what personal data the Kanban board application collects, why it is collected, how it is stored, and what rights users have over their data. This application is self-hosted on a private Raspberry Pi and is not a commercial service.

---

## Who is responsible for this data

This application is operated privately. The person who installed and administers the application (the administrator) is responsible for data stored on the system. If you have questions about your data, contact the administrator directly.

---

## What personal data we collect

### Account data

When you create an account, the following is stored:

| Data | Purpose |
|---|---|
| Your name | Displayed to other users, shown on task assignments and comments |
| Your email address | Used to log in and to send notifications |
| Your password | Stored as a bcrypt hash — the original password is never stored or readable |
| Account creation timestamp | System record |

We do not collect your address, phone number, date of birth, or any other personal information.

### Activity data

While you use the application, we store:

| Data | Purpose |
|---|---|
| Actions you take (moving tasks, creating tasks, commenting, etc.) | Activity log for audit trail and collaboration context |
| Timestamps of those actions | Part of the audit trail |
| Your IP address (in server logs) | Security — used to enforce rate limiting and detect abuse |

### Content data

Any content you create — task titles, descriptions, comments, and uploaded files — is stored and associated with your account. This content may be visible to other users who have access to the same boards.

---

## Legal basis for processing (GDPR)

Under the General Data Protection Regulation (GDPR), personal data must be processed on a lawful basis. For this application:

| Data type | Legal basis |
|---|---|
| Account data | Contractual necessity — required to provide you access to the service |
| Activity log | Legitimate interest — maintaining an audit trail for security and collaboration |
| Server logs (IP address) | Legitimate interest — security and abuse prevention |
| Content you create | Contractual necessity — the content is the core purpose of the service |

---

## How data is stored

All data is stored on a Raspberry Pi located on a private home network. Specifically:

- **Database** (user accounts, task data, activity logs): stored in PostgreSQL on a USB drive connected to the Pi
- **File attachments**: stored as files on the USB drive
- **Server logs**: stored temporarily in Docker container logs

The application uses HTTPS (TLS 1.2/1.3) to encrypt all data transmitted between your browser and the server. Passwords are hashed using bcrypt before storage — they cannot be recovered even by the administrator.

The Pi is not publicly accessible from the internet unless specifically configured to be. It is only accessible on the local network.

---

## Data retention

| Data | Retention period |
|---|---|
| Account data | Until you request deletion |
| Tasks, comments, attachments | Until deleted by you or a board owner |
| Activity log entries | 12 months, then automatically deleted |
| Server logs | 30 days, then automatically deleted |

---

## Who can see your data

| Audience | What they can see |
|---|---|
| Other users on the same board | Your name, tasks you create, comments you post, files you attach |
| Board owners | All activity on their boards |
| Administrator | All data stored on the system, including database contents |

Your password is never visible to anyone, including the administrator. It is stored only as a one-way hash.

---

## Your rights under GDPR

As a user, you have the following rights regarding your personal data:

### Right of access
You can request a copy of all personal data we hold about you. Contact the administrator and they will provide this within 30 days.

### Right to rectification
If any of your personal data is inaccurate, you can correct your name and email address yourself in your profile settings. For other corrections, contact the administrator.

### Right to erasure ("right to be forgotten")
You can request that your account and associated personal data be deleted. Note that:
- Content you created (tasks, comments) may be retained if other users depend on it, but it will be anonymised (your name removed)
- Activity log entries referencing your account will be anonymised
- Backups may retain your data until they expire or are deleted

To request deletion, contact the administrator.

### Right to data portability
You can request an export of your data in a machine-readable format (JSON or CSV). Contact the administrator.

### Right to restrict processing
You can request that we stop processing your personal data in certain circumstances. Contact the administrator.

### Right to object
You can object to the processing of your personal data where the legal basis is legitimate interest (activity logs, server logs). Contact the administrator.

---

## Cookies and tracking

This application uses a single session cookie to keep you logged in between visits. It does not use:

- Tracking cookies
- Analytics cookies
- Advertising cookies
- Third-party cookies of any kind

No data is shared with third parties. No external analytics services (Google Analytics, etc.) are used.

---

## Third-party services

This application does not send your data to any third-party services. All data stays on the Pi.

The only external network calls the application makes are:

- Email notifications, if configured — sent via an SMTP provider specified by the administrator
- Docker image pulls during updates — these are downloads from Docker Hub and do not transmit user data

---

## Data breaches

If the administrator becomes aware of a data breach (unauthorised access to the system), affected users will be notified as soon as reasonably possible. Because this is a private, locally-hosted application, the attack surface is limited to the local network unless the Pi has been exposed publicly.

---

## Changes to this policy

If this privacy policy changes significantly, users will be notified via the application. The date of the last update is recorded at the bottom of this document.

---

## Contact

For any questions or requests relating to your personal data, contact the administrator of this application directly.

---

*Last updated: April 2026*
