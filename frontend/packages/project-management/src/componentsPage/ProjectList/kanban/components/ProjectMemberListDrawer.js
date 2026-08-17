/**
 * 项目成员列表：搜索、改角色、移除；管理员/创建者可操作。
 */
'use client';

import { DeleteOutlined, PlusOutlined } from '@ant-design/icons';
import { Avatar, Button, Drawer, Empty, Input, Popconfirm, Select, Tag } from 'antd';
import { useTranslations } from 'next-intl';
import { PROJECT_MEMBER_ROLE_OPTIONS } from './constants';
import { getMemberInitial, getMemberLabel } from './helpers';

/**
 * @param {{
 *   open: boolean,
 *   project?: object|null,
 *   members?: Array,
 *   projectMembers?: Array,
 *   canManage?: boolean,
 *   search: string,
 *   onSearchChange: (value: string) => void,
 *   onClose: () => void,
 *   onInvite: () => void,
 *   onRemove: (userId: number) => void,
 *   onRoleChange: (userId: number, role: string) => void,
 * }} props
 */
export function ProjectMemberListDrawer({
	open,
	project,
	members = [],
	projectMembers = [],
	canManage = false,
	search,
	onSearchChange,
	onClose,
	onInvite,
	onRemove,
	onRoleChange,
}) {
	const t = useTranslations();

	const filteredMembers = projectMembers.filter((item) => {
		if (!search.trim()) return true;
		const member = members.find((row) => row.id === item.user_id);
		const label = getMemberLabel(member, item.user_id).toLowerCase();
		return label.includes(search.trim().toLowerCase());
	});

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
			width={420}
			open={open}
			onClose={onClose}
		>
			<Input.Search
				placeholder={t('kanban.memberSearchPlaceholder')}
				allowClear
				value={search}
				onChange={(event) => onSearchChange(event.target.value)}
				className="mb-4"
			/>
			<div className="flex flex-col">
				{filteredMembers.map((item) => {
					const member = members.find((row) => row.id === item.user_id);
					const name = getMemberLabel(member, item.user_id);
					const isOwner = item.user_id === project?.owner_id;
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
									{getMemberInitial(member, item.user_id)}
								</Avatar>
								<div className="min-w-0">
									<div className="font-medium text-sm truncate">{name}</div>
									<div className="text-xs text-base-label truncate">{member?.email}</div>
								</div>
							</div>
							<div className="flex items-center gap-1 shrink-0">
								{isOwner ? (
									<Tag color="purple" className="w-[88px] text-center m-0">
										{t('kanban.owner')}
									</Tag>
								) : (
									<Select
										size="small"
										value={item.role}
										className="w-[88px]"
										disabled={!canManage}
										onChange={(role) => onRoleChange(item.user_id, role)}
										options={PROJECT_MEMBER_ROLE_OPTIONS.map((option) => ({
											value: option.value,
											label: t(option.labelKey),
										}))}
									/>
								)}
								{canManage && !isOwner ? (
									<Popconfirm
										title={t('kanban.memberRemoveConfirm')}
										onConfirm={() => onRemove(item.user_id)}
										okText={t('common.confirm')}
										cancelText={t('common.cancel')}
									>
										<Button type="text" size="small" danger icon={<DeleteOutlined />} />
									</Popconfirm>
								) : null}
							</div>
						</div>
					);
				})}
				{filteredMembers.length === 0 ? (
					<Empty description={t('kanban.projectMembersEmpty')} className="mt-10" />
				) : null}
			</div>
		</Drawer>
	);
}
