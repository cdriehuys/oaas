import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import Message from './Message';

const client = new QueryClient();

function App() {

  return (
    <QueryClientProvider client={client}>
      <Message />
    </QueryClientProvider>
  )
}

export default App
