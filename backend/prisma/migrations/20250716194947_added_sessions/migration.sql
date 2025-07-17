/*
  Warnings:

  - You are about to drop the `post` table. If the table is not empty, all the data it contains will be lost.

*/
-- DropTable
DROP TABLE "post";

-- CreateTable
CREATE TABLE "active_sessions" (
    "id" SERIAL NOT NULL,
    "user_id" INTEGER NOT NULL,
    "session_uuid" TEXT NOT NULL,
    "last_logout" TIMESTAMP(3) NOT NULL,
    "exp" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "active_sessions_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "active_sessions_session_uuid_key" ON "active_sessions"("session_uuid");

-- CreateIndex
CREATE INDEX "active_sessions_user_id_session_uuid_idx" ON "active_sessions"("user_id", "session_uuid");

-- CreateIndex
CREATE INDEX "active_sessions_exp_idx" ON "active_sessions"("exp");

-- CreateIndex
CREATE UNIQUE INDEX "active_sessions_user_id_session_uuid_key" ON "active_sessions"("user_id", "session_uuid");

-- AddForeignKey
ALTER TABLE "active_sessions" ADD CONSTRAINT "active_sessions_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "User"("id") ON DELETE CASCADE ON UPDATE CASCADE;
