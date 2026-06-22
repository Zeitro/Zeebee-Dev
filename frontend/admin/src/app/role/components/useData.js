import { useState } from 'react';
import { createRequest } from '@repo/shared-methods/utils';

/**
 * fieldNames 将 value 映射为 id 时，TreeSelect 默认按 id 过滤；
 * 此函数改为按节点的 name / path 做不区分大小写的包含匹配。
 *
 * @param {string} input 搜索关键字
 * @param {{ name?: string, path?: string }} treeNode 树节点原始数据
 * @returns {boolean}
 */
export const filterTreeNodeByNameOrPath = (input, treeNode) => {
	const keyword = input.trim().toLowerCase();
	if (!keyword) return true;
	const name = String(treeNode.name ?? '').toLowerCase();
	const path = String(treeNode.path ?? '').toLowerCase();
	return name.includes(keyword) || path.includes(keyword);
};

const useData = () => {
	const { request } = createRequest();
	const [loading, setLoading] = useState(false);
	const [apiTreeData, setApiTreeData] = useState([]);

	const changeTreeLabel = (list) => {
		return list.map((item) => {
			item.label = item?.path ? (
				<>
					{item.name} (<span className="text-base-label"> {item.path}</span>)
				</>
			) : (
				item.name
			);
			item.value = item.id;
			if (item.children) {
				item.children = changeTreeLabel(item.children);
			}
			return item;
		});
	};

	const getApiTree = async (v) => {
		try {
			setLoading(true);
			const payload = new URLSearchParams({ type: v });
			const api = `/api/tree?${payload}`;
			const res = await request(api);
			// const apiList = changeTreeLabel(res?.list);
			setApiTreeData(res?.list);
			setLoading(false);
		} catch (error) {
			setLoading(false);
			throw error;
		}
	};
	return { getApiTree, apiTreeData, loading };
};

export default useData;
