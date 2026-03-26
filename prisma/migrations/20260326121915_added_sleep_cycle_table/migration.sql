-- CreateTable
CREATE TABLE "sleep_cycles" (
    "id" UUID NOT NULL,
    "user_id" UUID NOT NULL,
    "sleep_data" JSONB NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "sleep_cycles_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE INDEX "sleep_cycles_user_id_idx" ON "sleep_cycles"("user_id");

-- AddForeignKey
ALTER TABLE "sleep_cycles" ADD CONSTRAINT "sleep_cycles_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;
