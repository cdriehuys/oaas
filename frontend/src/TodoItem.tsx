import { useMutation } from "@tanstack/react-query";
import {
  putTodosByTodoIdState,
  type Todo,
  type TodoCollection,
} from "./client";
import { todoQueryKeys } from "./todoQueries";
import type React from "react";

interface Props {
  todo: Todo;
}

const TodoItem = ({ todo }: Props) => {
  const updateStateMutation = useMutation({
    mutationFn: async (isComplete: boolean) => {
      const state = isComplete ? "complete" : "incomplete";
      return putTodosByTodoIdState<true>({
        path: { todoID: todo.id },
        body: state,
      });
    },
    onMutate: async (isComplete, context) => {
      // Cancel outgoing queries
      await context.client.cancelQueries({ queryKey: todoQueryKeys.all });

      const completeKey = todoQueryKeys.listByState("complete");
      const incompleteKey = todoQueryKeys.listByState("incomplete");

      const previousCompleteItems: Todo[] =
        context.client.getQueryData<TodoCollection>(completeKey)?.items ?? [];
      const previousIncompleteItems: Todo[] =
        context.client.getQueryData<TodoCollection>(incompleteKey)?.items ?? [];

      if (isComplete) {
        const completedTodo: Todo = {
          ...todo,
          completedAt: new Date().toISOString(),
        };
        context.client.setQueryData(completeKey, {
          items: [completedTodo, ...previousCompleteItems],
        });
        context.client.setQueryData(incompleteKey, {
          items: previousIncompleteItems.filter((t) => t.id !== todo.id),
        });
      } else {
        const incompleteTodo: Todo = { ...todo };
        delete incompleteTodo["completedAt"];

        let insertionPoint = 0;
        for (const t of previousIncompleteItems) {
          if (t.id > todo.id) {
            break;
          }

          insertionPoint += 1;
        }

        const newIncompleteItems = [...previousIncompleteItems];
        newIncompleteItems.splice(insertionPoint, 0, incompleteTodo);

        context.client.setQueryData(completeKey, {
          items: previousCompleteItems.filter((t) => t.id !== todo.id),
        });
        context.client.setQueryData(incompleteKey, {
          items: newIncompleteItems,
        });
      }

      return () => {
        context.client.setQueryData(completeKey, {
          items: previousCompleteItems,
        });
        context.client.setQueryData(incompleteKey, {
          items: previousIncompleteItems,
        });
      };
    },

    onError: (_err, _vars, rollback) => rollback?.(),

    onSettled: (_data, _error, _variables, _onMutateResult, context) =>
      context.client.invalidateQueries({ queryKey: todoQueryKeys.all }),
  });

  const handleOnChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    updateStateMutation.mutate(e.currentTarget.checked);
  };

  let content = (
    <label>
      <input
        className="mr-2"
        type="checkbox"
        disabled={updateStateMutation.isPending}
        onChange={handleOnChange}
        checked={!!todo.completedAt}
      />
      {todo.title}
    </label>
  );

  if (todo.completedAt) {
    content = <span className="line-through">{content}</span>;
  }

  return content;
};

export default TodoItem;
