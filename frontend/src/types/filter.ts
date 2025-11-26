/**
 * Types for article filtering
 */

export interface FilterCondition {
  id: number;
  logic?: 'and' | 'or' | null;
  negate: boolean;
  field: string;
  operator?: string | null;
  value: string;
  values: string[];
}

export interface FieldOption {
  value: string;
  labelKey: string;
  multiSelect: boolean;
  booleanField?: boolean;
}

export interface OperatorOption {
  value: string;
  labelKey: string;
}

export interface LogicOption {
  value: 'and' | 'or';
  labelKey: string;
}
