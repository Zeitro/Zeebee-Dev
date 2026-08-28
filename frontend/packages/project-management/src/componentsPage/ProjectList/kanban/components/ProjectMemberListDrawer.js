/**
 * 项目成员列表：服务端搜索、改角色、移除；管理员/创建者可操作。
 * 展示列表与全量 projectMembers 分离，搜索仅影响本抽屉。
 */
'use client';

import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { DeleteOutlined, PlusOutlined } from '@ant-design/icons';
import { Avatar, Button, Drawer, Empty, Input, Popconfirm, Select, Tag } from 'antd';
import { useTranslations } from 'next-intl';
import { createRequest } from '@repo/shared-methods/utils';
import { createProjectMemberApi } from '@repo/project-management/api/projectMember';
import { localizeSelectOptions } from '@repo/project-management/enum/task';
import { getWorkflowMemberRoleI18nKey } from '@repo/project-management/enum/workflow';
import { useOptionsStore } from '@repo/project-management/store/options';
import { getMemberInitial } from './helpers';

/**
 * @param {{
 *   open: boolean,
 *   projectId?: number|string,
 *   projectMembers?: Array,
 *   canManage?: boolean,
 *   onClose: () => void,
 *   onInvite: () => void,
 *   onSearchMembers: (keyword: string) => Promise<Array>,
 *   onSuccess?: () => void | Promise<void>,
 *   messageApi?: import('antd').MessageInstance,
 * }} props
 */
export function ProjectMemberListDrawer({
	open,
	projectId,
	projectMembers = [],
	canManage = false,
	onClose,
	onInvite,
	onSearchMembers,
	onSuccess,
	messageApi,
}) {
	const t = useTranslations();
	const { request } = useMemo(() => createRequest(), []);
	const projectMemberApi = useMemo(() => createProjectMemberApi(request), [request]);
	const memberOptionsRaw = useOptionsStore((state) => state.memberOptions);
	const getOptions = useOptionsStore((state) => state.getOptions);
	const roleOptions = useMemo(
		() =>
			localizeSelectOptions(
				memberOptionsRaw.filter(
					(item) => item.value !== 'ProjectOwner' && item.value !== 'TaskCreator',
				),
				t,
				getWorkflowMemberRoleI18nKey,
			),
		[memberOptionsRaw, t],
	);

	const [search, setSearch] = useState('');
	const [activeKeyword, setActiveKeyword] = useState('');
	const [searchedResults, setSearchedResults] = useState([]);
	const [searchLoading, setSearchLoading] = useState(false);
	const [updatingMemberId, setUpdatingMemberId] = useState(null);
	const [removingMemberId, setRemovingMemberId] = useState(null);
	const activeKeywordRef = useRef('');

	const displayMembers = activeKeyword ? searchedResults : projectMembers;

	const resetSearch = useCallback(() => {
		setSearch('');
		setActiveKeyword('');
		activeKeywordRef.current = '';
		setSearchedResults([]);
	}, []);

	const runSearch = useCallback(
		async (keyword) => {
			setSearchLoading(true);
			try {
				const list = await onSearchMembers(keyword);
				setSearchedResults(list);
			} catch {
				setSearchedResults([]);
			} finally {
				setSearchLoading(false);
			}
		},
		[onSearchMembers],
	);

	const handleSearch = async (value) => {
		const keyword = value.trim();
		if (!keyword) {
			resetSearch();
			return;
		}
		setActiveKeyword(keyword);
		activeKeywordRef.current = keyword;
		await runSearch(keyword);
	};

	const handleSearchChange = (event) => {
		const value = event.target.value;
		setSearch(value);
		if (!value) resetSearch();
	};

	const handleClose = () => {
		resetSearch();
		onClose();
	};

	const handleRoleChange = useCallback(
		async (memberId, role) => {
			if (!projectId || memberId == null) return;
			setUpdatingMemberId(memberId);
			try {
				await projectMemberApi.updateRole(memberId, {
					project_id: Number(projectId),
					role,
				});
				messageApi.success(t('kanban.memberRoleUpdated'));
				await onSuccess?.();
			} catch (error) {
				messageApi.error(error?.message || t('kanban.memberRoleUpdateFailed'));
			} finally {
				setUpdatingMemberId(null);
			}
		},
		[projectId, projectMemberApi, onSuccess, t, messageApi],
	);

	const handleRemove = useCallback(
		async (memberId) => {
			if (!projectId || memberId == null) return;
			setRemovingMemberId(memberId);
			try {
				await projectMemberApi.delete(memberId, {
					project_id: Number(projectId),
				});
				messageApi.success(t('kanban.memberRemoved'));
				await onSuccess?.();
			} catch (error) {
				messageApi.error(error?.message || t('kanban.memberRemoveFailed'));
			} finally {
				setRemovingMemberId(null);
			}
		},
		[projectId, projectMemberApi, onSuccess, t, messageApi],
	);

	useEffect(() => {
		if (!open) return;
		// if (useOptionsStore.getState().memberOptions.length) return;
		getOptions('project_management_member_role', 'memberOptions').catch(() => {});
	}, [open, getOptions]);

	// 全量成员更新后（移除/改角色），搜索态下刷新抽屉列表
	useEffect(() => {
		if (!open || !activeKeywordRef.current) return undefined;

		let cancelled = false;
		(async () => {
			setSearchLoading(true);
			try {
				const list = await onSearchMembers(activeKeywordRef.current);
				if (!cancelled) setSearchedResults(list);
			} catch {
				if (!cancelled) setSearchedResults([]);
			} finally {
				if (!cancelled) setSearchLoading(false);
			}
		})();

		return () => {
			cancelled = true;
		};
	}, [open, projectMembers, onSearchMembers]);

	return (
		<Drawer
			title={
				<div className="flex items-center justify-between">
					<span>{t('kanban.projectMembers')}</span>
					{canManage ? (
						<Button type="text" size="small" icon={<PlusOutlined />} onClick={onInvite}>
							{t('kanban.memberAdd')}
						</Button>
					) : null}
				</div>
			}
			placement="right"
			width={500}
			open={open}
			onClose={handleClose}
		>
			<Input.Search
				placeholder={t('kanban.memberSearchPlaceholder')}
				allowClear
				value={search}
				loading={searchLoading}
				onChange={handleSearchChange}
				onSearch={handleSearch}
				className="mb-4"
			/>
			<div className="flex flex-col">
				{displayMembers.map((item) => {
					const isOwner = item.role === 'ProjectOwner';
					return (
						<div
							key={item.id || item.user_id}
							className="flex items-center justify-between py-3 border-b border-base-border"
						>
							<div className="flex items-center gap-3 min-w-0">
								<Avatar
									size={36}
									className="shrink-0"
									style={{
										backgroundColor: isOwner
											? 'var(--purple-dark)'
											: 'var(--blue-medium)',
									}}
								>
									{getMemberInitial({
										full_name: item.full_name,
										email: item.email,
									})}
								</Avatar>
								<div className="min-w-0">
									<div className="font-medium text-sm truncate">
										{item.full_name}
									</div>
									<div className="text-xs text-base-label truncate">
										{item.email}
									</div>
								</div>
							</div>
							<div className="flex items-center gap-1 shrink-0">
								{isOwner ? (
									<Tag
										color="purple"
										className="w-[88px] text-center m-0 !text-sm"
									>
										{t('common.roleProjectOwner')}
									</Tag>
								) : (
									<Select
										size="small"
										value={item.role}
										className="w-[110px]"
										disabled={!canManage}
										loading={updatingMemberId === item.id}
										onChange={(role) => handleRoleChange(item.id, role)}
										options={roleOptions}
									/>
								)}
								{canManage && !isOwner ? (
									<Popconfirm
										title={t('kanban.memberRemoveConfirm')}
										onConfirm={() => handleRemove(item.id)}
										okText={t('common.confirm')}
										cancelText={t('common.cancel')}
									>
										<Button
											type="text"
											size="small"
											danger
											loading={removingMemberId === item.id}
											icon={<DeleteOutlined />}
										/>
									</Popconfirm>
								) : (
									<span className="w-[14px] inline-block"></span>
								)}
							</div>
						</div>
					);
				})}
				{!searchLoading && displayMembers.length === 0 ? (
					<Empty description={t('kanban.projectMembersEmpty')} className="mt-10" />
				) : null}
			</div>
		</Drawer>
	);
}
