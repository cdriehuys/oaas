import type React from "react";
import { useTodos } from "./todoQueries";
import TodoItem from "./TodoItem";

const TodoWrapper = ({ children }: { children?: React.ReactNode }) => {
  return <div className="mt-16 mx-auto max-w-2xl">{children}</div>;
};

interface Props {
  state: "complete" | "incomplete";
  title: string;
}

export default function TodoList({ state, title }: Props) {
  const query = useTodos(state);

  if (query.isLoading) {
    return <p>Loading...</p>;
  }

  if (query.error) {
    return <p>Failed to fetch todos: {query.error.toString()}</p>;
  }

  const todos = query.data?.data?.items ?? [];

  return (
    <TodoWrapper>
      <h1 className="mb-8 text-4xl text-center text-zinc-700 font-light italic">
        {title}
      </h1>
      {todos.length === 0 ? (
        <p>Nothing!</p>
      ) : (
        <ul>
          {todos.map((t) => (
            <li key={t.id}>
              <TodoItem todo={t} />
            </li>
          ))}
        </ul>
      )}
    </TodoWrapper>
  );
}
