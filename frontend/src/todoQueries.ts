import { useQuery } from "@tanstack/react-query";
import { getTodos } from "./client";

export const todoQueryKeys = {
    all: ['todos'] as const,
    list: () => [todoQueryKeys.all, 'list'] as const,
};

export const useTodos = () => useQuery({
    queryKey: todoQueryKeys.list(),
    queryFn: () => getTodos<true>()
});
