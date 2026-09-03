/**
 * 邀请项目成员：candidate 拉候选、addMembers 批量添加；两处入口共用本弹窗。
 */
'use client';

import { useCallback, useEffect, useMemo, useState } from 'react';
import { Avatar, Checkbox, Empty, Input, message } from 'antd';
import { useTranslations } from 'next-intl';
import { NewModal } from '@repo/shared-methods';
import { createRequest } from '@repo/shared-methods/utils';
import {
	createProjectMemberApi,
	PROJECT_MEMBER_CANDIDATE_PURPOSE_ADD_MEMBER,
	PROJECT_MEMBER_DEFAULT_ROLE,
} from '@repo/project-management/api/projectMember';
import { mapProjectMember } from '@repo/project-management/lib/mapProjectMember';
import { getMemberInitial, getMemberLabel, unwrapList } from './helpers';

/**
 * @param {object|null|undefined} item
 * @returns {number|undefined}
 */
function getCandidateUserId(item) {
	const id = item?.organization_user_id ?? item?.user_id ?? item?.id;
	return id == null ? undefined : Number(id);
}

/**
 * @param {{
 *   open: boolean,
 *   projectId?: number|string,
 *   onCancel: () => void,
 *   onSuccess?: () => void | Promise<void>,
 * }} props
 */
export function InviteMemberModal({ open, projectId, onCancel, onSuccess, messageApi }) {
	const t = useTranslations();
	const { request } = useMemo(() => createRequest(), []);
	const projectMemberApi = useMemo(() => createProjectMemberApi(request), [request]);

	const [search, setSearch] = useState('');
	const [selected, setSelected] = useState(new Set());
	const [candidates, setCandidates] = useState([]);
	const [candidatesLoading, setCandidatesLoading] = useState(false);
	const [inviteLoading, setInviteLoading] = useState(false);

	const resetState = useCallback(() => {
		setSearch('');
		setSelected(new Set());
		setCandidates([]);
	}, []);

	const loadCandidates = useCallback(
		async (keyword = '') => {
			if (!projectId) {
				setCandidates([]);
				return;
			}
			setCandidatesLoading(true);
			try {
				const res = await projectMemberApi.candidate({
					project_id: Number(projectId),
					keyword: keyword.trim(),
					purpose: PROJECT_MEMBER_CANDIDATE_PURPOSE_ADD_MEMBER,
				});
				setCandidates(unwrapList(res).map(mapProjectMember).filter(Boolean));
			} catch {
				setCandidates([]);
			} finally {
				setCandidatesLoading(false);
			}
		},
		[projectId, projectMemberApi],
	);

	useEffect(() => {
		if (!open) return undefined;

		let cancelled = false;
		(async () => {
			if (!projectId) return;
			setCandidatesLoading(true);
			try {
				const res = await projectMemberApi.candidate({
					project_id: Number(projectId),
					keyword: '',
					purpose: PROJECT_MEMBER_CANDIDATE_PURPOSE_ADD_MEMBER,
				});
				if (!cancelled) {
					setCandidates(unwrapList(res).map(mapProjectMember).filter(Boolean));
				}
			} catch {
				if (!cancelled) setCandidates([]);
			} finally {
				if (!cancelled) setCandidatesLoading(false);
			}
		})();

		return () => {
			cancelled = true;
		};
	}, [open, projectId, projectMemberApi]);

	const handleSearch = async (value) => {
		await loadCandidates(value.trim());
	};

	const handleSearchChange = (event) => {
		const value = event.target.value;
		setSearch(value);
		if (!value) loadCandidates('');
	};

	const toggleSelection = (userId) => {
		setSelected((prev) => {
			const next = new Set(prev);
			if (next.has(userId)) next.delete(userId);
			else next.add(userId);
			return next;
		});
	};

	const handleClose = () => {
		resetState();
		onCancel();
	};

	const handleInvite = async () => {
		const userIds = [...selected];
		if (!projectId || !userIds.length) return;

		setInviteLoading(true);
		try {
			await projectMemberApi.addMembers({
				project_id: Number(projectId),
				members: userIds.map((organization_user_id) => ({
					organization_user_id: Number(organization_user_id),
					role: PROJECT_MEMBER_DEFAULT_ROLE,
				})),
			});
			messageApi.success(t('kanban.inviteSuccess', { count: userIds.length }));
			await onSuccess?.();
			resetState();
			onCancel();
		} catch (error) {
			messageApi.error(error?.message || t('kanban.inviteFailed', { count: userIds.length }));
			throw error;
		} finally {
			setInviteLoading(false);
		}
	};

	if (!open) return null;

	return (
		<NewModal
			title={t('kanban.inviteMembers')}
			width={500}
			destroyOnHidden
			onCancel={handleClose}
			onOk={handleInvite}
			okText={t('kanban.inviteAdd')}
			okButtonProps={{
				disabled: selected.size === 0,
			}}
			confirmLoading={inviteLoading}
			cancelText={t('common.cancel')}
		>
			<Input.Search
				placeholder={t('kanban.inviteSearchPlaceholder')}
				allowClear
				value={search}
				loading={candidatesLoading}
				onChange={handleSearchChange}
				onSearch={handleSearch}
				className="mb-4"
			/>
			<div className="min-h-[360px] max-h-[600px] overflow-y-auto">
				{!candidatesLoading && candidates.length === 0 ? (
					<Empty
						description={t('kanban.inviteEmpty')}
						image={Empty.PRESENTED_IMAGE_SIMPLE}
					/>
				) : (
					candidates.map((item) => {
						const userId = getCandidateUserId(item);
						if (userId == null) return null;
						return (
							<div
								key={userId}
								className={`flex items-center justify-between px-3 py-2.5 border-b border-base-border cursor-pointer ${
									selected.has(userId) ? 'bg-blue-light-1' : ''
								}`}
								onClick={() => toggleSelection(userId)}
							>
								<div className="flex items-center gap-2.5">
									<Checkbox
										checked={selected.has(userId)}
										onChange={() => toggleSelection(userId)}
										onClick={(event) => event.stopPropagation()}
									/>
									<Avatar
										size={36}
										style={{ backgroundColor: 'var(--blue-medium)' }}
										className="shrink-0"
									>
										{getMemberInitial(item, userId)}
									</Avatar>
									<div>
										<div className="font-medium">
											{getMemberLabel(item, userId)}
										</div>
										<div className="text-xs text-base-label">{item.email}</div>
									</div>
								</div>
							</div>
						);
					})
				)}
			</div>
		</NewModal>
	);
}
