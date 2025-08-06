import { z } from "zod";

export const PageInfoSchema = z.object({
  currentPage: z.number(),
  totalPages: z.number(),
  hasNextPage: z.boolean(),
  hasPrevPage: z.boolean(),
});

export function createPaginationSchema<T extends z.ZodTypeAny>(itemSchema: T) {
  return z.object({
    pageInfo: PageInfoSchema,
    items: z.array(itemSchema),
  });
}
