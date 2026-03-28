-- CreateTable
CREATE TABLE "workout_splits" (
    "id" UUID NOT NULL,
    "trainer_id" UUID NOT NULL,
    "user_id" UUID NOT NULL,
    "split_data" JSONB NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "workout_splits_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "workout_splits_user_id_key" ON "workout_splits"("user_id");

-- AddForeignKey
ALTER TABLE "workout_splits" ADD CONSTRAINT "workout_splits_trainer_id_fkey" FOREIGN KEY ("trainer_id") REFERENCES "trainers"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "workout_splits" ADD CONSTRAINT "workout_splits_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;
