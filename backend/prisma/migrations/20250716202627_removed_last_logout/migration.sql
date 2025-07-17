/*
  Warnings:

  - You are about to drop the column `last_logout` on the `active_sessions` table. All the data in the column will be lost.

*/
-- AlterTable
ALTER TABLE "active_sessions" DROP COLUMN "last_logout";
