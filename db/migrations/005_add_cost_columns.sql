-- Migration: Add downtime_loss to repair_cost_details
ALTER TABLE repair_cost_details ADD COLUMN downtime_loss DECIMAL(12,2) DEFAULT 0;
