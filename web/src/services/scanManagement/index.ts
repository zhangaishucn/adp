import request from '../request';
import * as ScanManagement from './type';

const BASE_URL = '/api/data-connection/v1/metadata';
const EXCEL_BASE_URL = '/api/data-connection/v1/gateway/excel';

export const createScanTask = (params: ScanManagement.ScanRequest) => {
  const obj = {
    ...params,
    use_default_template: true,
    field_list_when_change: [],
    use_multi_threads: true,
    tables: params.tables || [],
  };
  return request.post<ScanManagement.ScanResponse>(`${BASE_URL}/scan`, obj);
};

export const batchCreateScanTask = (params: ScanManagement.ScanRequest[]) => {
  const obj = params.map((item) => ({
    ...item,
    use_default_template: true,
    field_list_when_change: [],
    use_multi_threads: true,
    tables: item.tables || [],
  }));
  return request.post<ScanManagement.ScanResponse>(`${BASE_URL}/scan/batch`, obj);
};

export const getScanTaskStatus = (taskId: string) => {
  return request.get<ScanManagement.TaskStatusResponse>(`${BASE_URL}/scan/status/${taskId}`);
};

export const getScanTaskTableStatus = (params: { id: string; type: string[] }) => {
  return request.post<ScanManagement.TableStatusResponse>(`${BASE_URL}/scan/status`, params);
};

export const retryScanTask = (params: ScanManagement.RetryRequest) => {
  return request.post<ScanManagement.RetryResponse>(`${BASE_URL}/scan/retry`, params);
};

export const getScanTaskList = (params: ScanManagement.PageQueryParams) => {
  const curSort = params.sort || 'start_time';
  const queryParams = {
    sort: curSort,
    direction: params.direction || 'desc',
    limit: params.limit || -1,
    keyword: params.keyword,
    offset: params.offset || 0,
    ds_id: params.ds_id,
    status: !params.status || params.status === 'all' ? '' : params.status,
    ...params.filters,
  };
  return request.get<ScanManagement.ScanTaskListResponse>(`${BASE_URL}/scan`, queryParams);
};

export const getDataSourceTables = (id: string, params?: ScanManagement.PageQueryParams) => {
  return request.get<ScanManagement.TableListResponse>(`${BASE_URL}/data-source/${id}`, params);
};

export const getScanTaskInfo = (taskId: string, params?: ScanManagement.ScanTaskInfoParams) => {
  return request.get<ScanManagement.ScanTaskInfoResponse>(`${BASE_URL}/scan/info/${taskId}`, params);
};

export const getTableColumns = (tableId: string, params?: ScanManagement.PageQueryParams) => {
  return request.get<ScanManagement.ColumnListResponse>(`${BASE_URL}/table/${tableId}`, params);
};

export const createExcelTable = (params: ScanManagement.CreateExcelTableRequest) => {
  return request.post<ScanManagement.ExcelTableResponse>(`${EXCEL_BASE_URL}/table`, params);
};

export const deleteExcelTable = (tableId: string) => {
  return request.delete<ScanManagement.ExcelTableResponse>(`${EXCEL_BASE_URL}/table/${tableId}`);
};

export const getExcelColumns = (params: ScanManagement.GetExcelColumnsRequest) => {
  return request.post<ScanManagement.ExcelColumnsResponse>(`${EXCEL_BASE_URL}/columns`, params, { timeout: 60000 });
};

export const getExcelSheets = (catalog: string, fileName: string) => {
  return request.get<ScanManagement.ExcelSheetListResponse>(
    `${EXCEL_BASE_URL}/sheet`,
    {
      catalog,
      file_name: fileName,
    },
    { timeout: 60000 }
  );
};

export const getExcelFiles = (catalog: string) => {
  return request.get<ScanManagement.ExcelFileListResponse>(`${EXCEL_BASE_URL}/files/${catalog}`, {}, { timeout: 60000 });
};

export default {
  createScanTask,
  batchCreateScanTask,
  getScanTaskStatus,
  getScanTaskTableStatus,
  retryScanTask,
  getScanTaskList,
  getDataSourceTables,
  getTableColumns,
  createExcelTable,
  deleteExcelTable,
  getExcelColumns,
  getExcelSheets,
  getExcelFiles,
};
