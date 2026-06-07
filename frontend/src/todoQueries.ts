import { useQuery } from "@tanstack/react-query";
import { getTodos } from "./client";

export const todoQueryKeys = {
  all: ["todos"] as const,
  lists: () => [...todoQueryKeys.all, "lists"] as const,
  listByState: (state: string) => [...todoQueryKeys.lists(), state] as const,
};

export const useTodos = (state: "complete" | "incomplete") =>
  useQuery({
    queryKey: todoQueryKeys.listByState(state),
    queryFn: () => getTodos<true>({ query: { state } }),
  });
