import Request from '../request';
import * as ActionType from './type';

const BASE_URL = '/api/ontology-manager/v1/knowledge-networks';

/**
 * 获取行动类列表
 * @param knId 知识网络ID
 * @param params 查询参数
 */
export const getActionTypes = (knId: string, params: ActionType.GetActionTypesRequest): Promise<ActionType.GetActionTypesResponse> => {
  return Request.get(`${BASE_URL}/${knId}/action-types`, params);
};

/**
 * 删除行动类
 * @param knId 知识网络ID
 * @param atIds 行动类ID列表
 */
export const deleteActionType = (knId: string, atIds: string[]): Promise<void> => {
  return Request.delete(`${BASE_URL}/${knId}/action-types/${atIds.join(',')}`);
};

/**
 * 新建行动类
 * @param knId 知识网络ID
 * @param data 创建数据
 */
export const createActionType = (knId: string, data: ActionType.CreateActionTypeRequest): Promise<void> => {
  return Request.post(`${BASE_URL}/${knId}/action-types`, { entries: data });
};

/**
 * 编辑行动类
 * @param knId 知识网络ID
 * @param atId 行动类ID
 * @param data 编辑数据
 */
export const editActionType = (knId: string, atId: string, data: ActionType.EditActionTypeRequest): Promise<void> => {
  return Request.put(`${BASE_URL}/${knId}/action-types/${atId}`, data);
};

/**
 * 获取行动类详情
 * @param knId 知识网络ID
 * @param atIds 行动类ID列表
 */
export const getActionTypeDetail = (knId: string, atIds: string[]): Promise<ActionType.ActionType[]> => {
  return Request.get<{ entries: ActionType.ActionType[] }>(`${BASE_URL}/${knId}/action-types/${atIds.join(',')}`).then((response) => response.entries);
};

export default {
  getActionTypes,
  deleteActionType,
  createActionType,
  editActionType,
  getActionTypeDetail,
};
