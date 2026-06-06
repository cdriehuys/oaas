import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import TodoList from './TodoList';

const client = new QueryClient();

function App() {

  return (
    <QueryClientProvider client={client}>
      <TodoList />
    </QueryClientProvider>
  )
}

export default App
