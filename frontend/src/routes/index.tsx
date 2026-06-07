import { createFileRoute } from '@tanstack/react-router'
import TodoList from '../TodoList'

export const Route = createFileRoute('/')({
  component: Index,
})

function Index() {
  return (
    <>
          <TodoList state="incomplete" title='Todo List' />
          <TodoList state="complete" title="Done" />
          </>
  )
}
