import { type Ref } from 'vue';
import type { ActionOption } from './useRuleOptions';

export function useRuleActions(actionOptions: ActionOption[]) {
  function addAction(actions: Ref<string[]>): void {
    const selectedActions = new Set(actions.value);
    const available = actionOptions.find((opt) => !selectedActions.has(opt.value));
    if (available) {
      actions.value.push(available.value);
    }
  }

  function removeAction(actions: Ref<string[]>, index: number): void {
    actions.value.splice(index, 1);
  }

  function updateAction(actions: Ref<string[]>, index: number, value: string): void {
    actions.value[index] = value;
  }

  function getAvailableActions(actions: Ref<string[]>, currentValue: string): ActionOption[] {
    const selectedActions = new Set(actions.value);
    return actionOptions.filter(
      (opt) => !selectedActions.has(opt.value) || opt.value === currentValue
    );
  }

  return {
    addAction,
    removeAction,
    updateAction,
    getAvailableActions,
  };
}
