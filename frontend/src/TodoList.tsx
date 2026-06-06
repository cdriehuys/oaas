import { useTodos } from "./todoQueries";

export default function TodoList() {
    const query = useTodos();

    if (query.isLoading) {
        return <p>Loading...</p>
    }

    if (query.error) {
        return <p>Failed to fetch todos: {query.error.toString()}</p>
    }

    const todos = query.data?.data?.items ?? [];

    return <>
        <h1>To-do List:</h1>
        {todos.length === 0 ? <p>Nothing!</p> : (
            <ul>
                {todos.map(t => <li key={t.id}>{t.title}</li>)}
            </ul>
        )}
    </>
}
