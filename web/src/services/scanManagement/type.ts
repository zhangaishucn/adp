export type Status = 'wait' | 'running' | 'success' | 'fail' | 'all';

export interface ScanRequest {
  scan_name: string;
  type: number;
  ds_info: { ds_id: string; ds_type: string };
  use_default_template?: boolean;
  field_list_when_change?: string[];
  use_multi_threads?: boolean;
  cron_expression?: string;
  tables?: string[];
}

export interface ScanResponse {
  id: string;
  ds_id: string;
  status: Status;
}

export interface TaskProcessInfo {
  table_count: number;
  success_count: number;
  fail_count: number;
}

export interface TaskResultInfo {
  table_count: number;
  success_count: number;
  fail_count: number;
  fail_stage: number;
  error_stack: string;
}

export interface TaskStatusResponse {
  id: string;
  status: Status;
  data: {
    task_process_info: TaskProcessInfo;
    task_result_info: TaskResultInfo;
  };
}

export interface TableStatus {
  table_id: string;
  table_name: string;
  status: Status;
  start_time: string;
}

export interface TableStatusResponse {
  id: string;
  tables: TableStatus[];
}

export interface RetryRequest {
  id: string;
  tables: string[];
}

export interface RetryResponse {
  id: string;
  tables: Array<{
    table_id: string;
    status: Status;
  }>;
}

export interface ScanTaskItem {
  id: string;
  name: string;
  create_user: string;
  scan_status: Status;
  start_time: string;
  type: number;
  ds_type: string;
  allow_multi_table_scan: boolean;
  task_process_info: string;
  task_result_info: string;
}

export interface ScanTaskListResponse {
  total_count: number;
  entries: ScanTaskItem[];
}

export interface TableInfo {
  id: string;
  name: string;
  create_time?: string;
  advanced_params: GetExcelColumnsRequest;
}

export interface TableListResponse {
  total_count: number;
  entries: TableInfo[];
}

export interface ColumnInfo {
  id: string;
  name: string;
  comment: string;
  type: string;
}

export interface ColumnListResponse {
  total_count: number;
  entries: ColumnInfo[];
}

export interface PageQueryParams {
  limit?: number;
  offset?: number;
  keyword?: string;
  sort?: string;
  direction?: string;
  name?: string;
  ds_id?: string;
  status?: Status;
  filters?: Record<string, any>;
}

export interface ScanTaskInfoParams {
  limit?: number;
  offset?: number;
  keyword?: string;
  sort?: string;
  direction?: string;
  status?: Status;
}

export interface ScanTaskInfoResponseItem {
  task_id: string;
  table_id: string;
  table_name: string;
  scan_status: Status;
  start_time: string;
}

export interface ScanTaskInfoResponse {
  total_count: number;
  entries: ScanTaskInfoResponseItem[];
}

export interface CreateExcelTableRequest {
  datasource_id: string;
  catalog: string;
  file_name: string;
  table_name: string;
  sheet?: string;
  all_sheet?: boolean;
  sheet_as_new_column?: boolean;
  start_cell: string;
  end_cell: string;
  has_headers?: boolean;
  columns: Array<{
    column: string;
    type: string;
  }>;
}

export interface ExcelTableResponse {
  tableId: string;
  tableName: string;
}

export interface Column {
  column: string;
  type: string;
}

export interface ExcelColumnsResponse {
  data: Column[];
}

export interface GetExcelColumnsRequest {
  id?: string;
  create_time?: string;
  datasource_id?: string;
  catalog: string;
  file_name: string;
  sheet?: string;
  all_sheet?: boolean;
  sheet_as_new_column?: boolean;
  start_cell: string;
  end_cell: string;
  has_headers?: boolean;
}

export interface ExcelFileListResponse {
  data: string[];
  total: number;
}

export interface ExcelSheetListResponse {
  data: string[];
  total: number;
}
