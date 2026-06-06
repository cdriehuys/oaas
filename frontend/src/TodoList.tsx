import type React from "react";
import { useTodos } from "./todoQueries";
import type { Todo } from "./client";

const TodoWrapper = ({ children }: { children?: React.ReactNode }) => {
    return <div className="mt-16 mx-auto max-w-2xl">{children}</div>
}

const TodoItem = ({ todo }: { todo: Todo}) => {
    return (
        <label>
            <input className="mr-2" type="checkbox" />
            {todo.title}
        </label>
    )
}

export default function TodoList() {
    const query = useTodos();

    if (query.isLoading) {
        return <p>Loading...</p>
    }

    if (query.error) {
        return <p>Failed to fetch todos: {query.error.toString()}</p>
    }

    const todos = query.data?.data?.items ?? [];

    return <TodoWrapper>
        <h1 className="mb-8 text-4xl text-center text-zinc-700 font-light italic">Todo List</h1>
        {todos.length === 0 ? <p>Nothing!</p> : (
            <ul>
                {todos.map(t => <li key={t.id}><TodoItem todo={t} /></li>)}
            </ul>
        )}
    </TodoWrapper>
}
