/**
 * 多执行者芯片：胶囊 chip + 右上角圆形删除；右侧圆形「+」打开成员菜单。
 */
'use client';

import { CheckOutlined, CloseOutlined, PlusOutlined } from '@ant-design/icons';
import { Avatar, Button, Dropdown } from 'antd';
import { useCallback, useEffect, useMemo, useState } from 'react';
import {
	PROJECT_MEMBER_CANDIDATE_PURPOSE_PROJECT_MEMBER,
	createProjectMemberApi,
} from '@repo/project-management/api/projectMember';
import { mapProjectMember } from '@repo/project-management/lib/mapProjectMember';
import { createRequest, displayValue } from '@repo/shared-methods/utils';
import { getMemberInitial, getMemberLabel, unwrapList } from '../helpers';

/**
 * candidate 成员 → AssigneeChips 内部结构（id 对齐 assignee_ids）。
 * @param {object} item
 * @returns {{ id: number, name?: string, email?: string, full_name?: string }|null}
 */
function mapAssigneeCandidate(item) {
	if (!item) return null;
	const id = item.organization_user_id ?? item.user_id;
	if (id == null) return null;
	return {
		id,
		name: item.full_name || item.name,
		email: item.email,
		full_name: item.full_name,
	};
}

/**
 * @param {{
 *   value?: number[],
 *   onChange?: (ids: number[]) => void,
 *   options?: Array,
 *   fallbacks?: Array,
 *   projectId?: number|string,
 * }} props
 */
export function AssigneeChips({ value, onChange, projectId }) {
	const ids = value || [];
	const { request } = useMemo(() => createRequest(), []);
	const projectMemberApi = useMemo(() => createProjectMemberApi(request), [request]);
	const [options, setOptions] = useState([]);

	const fetchAssigneeOptions = useCallback(async () => {
		if (!projectId) return [];
		try {
			const res = await projectMemberApi.candidate({
				project_id: Number(projectId),
				keyword: '',
				purpose: PROJECT_MEMBER_CANDIDATE_PURPOSE_PROJECT_MEMBER,
			});
			return unwrapList(res).map(mapProjectMember).map(mapAssigneeCandidate).filter(Boolean);
		} catch {
			return [];
		}
	}, [projectId, projectMemberApi]);

	useEffect(() => {
		let cancelled = false;
		(async () => {
			const nextOptions = await fetchAssigneeOptions();
			if (cancelled) return;
			setOptions(nextOptions);
		})();
		return () => {
			cancelled = true;
		};
	}, [fetchAssigneeOptions]);

	const findMember = useCallback((id) => options.find((member) => member.id === id), [options]);

	const menuItems = options.map((member) => {
		const name = getMemberLabel(member);
		const selected = ids.includes(member.id);
		return {
			key: String(member.id),
			label: (
				<span className="inline-flex items-center gap-2.5 min-w-[160px] py-0.5">
					<Avatar
						size={22}
						style={{ backgroundColor: 'var(--blue-medium)', fontSize: 11 }}
					>
						{getMemberInitial(member)}
					</Avatar>
					<span className="flex-1 truncate text-[13px] text-base-text">
						{displayValue(name)}
					</span>
					{selected ? (
						<CheckOutlined className="text-base-primary shrink-0 text-[12px]" />
					) : (
						<span className="w-3 shrink-0" />
					)}
				</span>
			),
			onClick: () =>
				onChange?.(selected ? ids.filter((id) => id !== member.id) : [...ids, member.id]),
		};
	});

	return (
		<div className="flex items-center gap-2 flex-wrap">
			{ids.map((id) => {
				const member = findMember(id);
				const name = getMemberLabel(member, id);
				return (
					<span key={id} className="relative inline-flex mt-1 mr-1">
						<span className="inline-flex items-center gap-1.5 rounded-full bg-base-text2 py-[3px] pl-[3px] pr-2.5">
							<Avatar
								size={22}
								style={{ backgroundColor: 'var(--blue-medium)', fontSize: 11 }}
							>
								{getMemberInitial(member, id)}
							</Avatar>
							<span className="text-[13px] text-base-text leading-none max-w-[100px] truncate">
								{displayValue(name)}
							</span>
						</span>
						<button
							type="button"
							aria-label="remove"
							className="absolute -top-1 -right-1 z-[1] flex h-[14px] w-[14px] items-center justify-center rounded-full border border-base-border bg-white text-base-label shadow-sm hover:border-base-error hover:text-base-error cursor-pointer p-0 leading-none"
							onClick={(event) => {
								event.stopPropagation();
								onChange?.(ids.filter((item) => item !== id));
							}}
						>
							<CloseOutlined className="text-[8px]" />
						</button>
					</span>
				);
			})}

			<Dropdown
				menu={{ items: menuItems }}
				trigger={['click']}
				placement="bottomLeft"
				disabled={menuItems.length === 0}
			>
				<Button
					type="text"
					shape="circle"
					size="small"
					disabled={menuItems.length === 0}
					className="!w-7 !h-7 !min-w-7 !p-0 bg-base-text2 hover:!bg-base-border/40 text-base-gray shrink-0"
					icon={<PlusOutlined className="text-xs" />}
				/>
			</Dropdown>
		</div>
	);
}
