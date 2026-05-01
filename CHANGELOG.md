# Changelog

## [1.0.0] - 2026-04-30

### 🚀 Milestone: Initial Release
This is the first stable release of the EMS Equipment Management System.

### Added
- **Equipment Management**: Full lifecycle tracking, category management, and workshop/factory organization.
- **QR Code Support**: Real-time generation and download of equipment QR codes (scans as equipment ID).
- **Core Operations**: Inspection, Maintenance plans, Repair ticketing, and Spare parts inventory.
- **AI Agent Integration**: Intelligent maintenance recommendations, failure analysis, and multi-bot Feishu (Lark) integration.
- **Analytics**: MTTR/MTBF tracking and failure trend analysis.
- **Security**: JWT authentication with role-based access control (Admin/Engineer/Operator/Maintenance).
- **Deployment**: Comprehensive Docker Compose setup with Nginx reverse proxy.

### Fixed
- Backend route definitions and corrupted setup functions.
- Frontend variable redeclaration in Equipment list view.
- Feishu Webhook verification sequencing instructions.
