import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import TodoList from './TodoList';

const client = new QueryClient();

function App() {

  return (
    <QueryClientProvider client={client}>
      <TodoList state="incomplete" title='Todo List' />
      <TodoList state="complete" title="Done" />
    </QueryClientProvider>
  )
}

export default App
