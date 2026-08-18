/**
 * 邀请项目成员：从未加入项目的租户成员中多选添加。
 */
'use client';

import { useMemo, useState } from 'react';
import { Avatar, Button, Checkbox, Empty, Input } from 'antd';
import { useTranslations } from 'next-intl';
import { NewModal } from '@repo/shared-methods';
import { getMemberInitial, getMemberLabel } from './helpers';

/**
 * @param {{
 *   open: boolean,
 *   members?: Array,
 *   projectMembers?: Array,
 *   loading?: boolean,
 *   onCancel: () => void,
 *   onInvite: (userIds: number[]) => Promise<void>,
 * }} props
 */
export function InviteMemberModal({
	open,
	members = [],
	projectMembers = [],
	loading = false,
	onCancel,
	onInvite,
}) {
	const t = useTranslations();
	const [search, setSearch] = useState('');
	const [selected, setSelected] = useState(new Set());

	const availableMembers = useMemo(() => {
		const projectMemberIds = new Set(projectMembers.map((item) => item.user_id));
		return members.filter((item) => !projectMemberIds.has(item.id));
	}, [members, projectMembers]);

	const filteredMembers = useMemo(() => {
		if (!search.trim()) return availableMembers;
		const keyword = search.trim().toLowerCase();
		return availableMembers.filter((item) => {
			const label = getMemberLabel(item).toLowerCase();
			return label.includes(keyword) || (item.email || '').toLowerCase().includes(keyword);
		});
	}, [availableMembers, search]);

	const toggleSelection = (userId) => {
		setSelected((prev) => {
			const next = new Set(prev);
			if (next.has(userId)) next.delete(userId);
			else next.add(userId);
			return next;
		});
	};

	const handleClose = () => {
		setSearch('');
		setSelected(new Set());
		onCancel();
	};

	const handleInvite = async () => {
		await onInvite([...selected]);
		setSelected(new Set());
		setSearch('');
	};

	if (!open) return null;

	return (
		<NewModal
			title={t('kanban.inviteMembers')}
			width={480}
			destroyOnHidden
			onCancel={handleClose}
			footer={
				<div className="flex items-center justify-between">
					<span className="text-[13px] text-base-label">
						{selected.size > 0
							? t('kanban.inviteSelected', { count: selected.size })
							: ''}
					</span>
					<Button
						type="primary"
						loading={loading}
						disabled={selected.size === 0}
						onClick={handleInvite}
					>
						{t('kanban.inviteAdd')}
					</Button>
				</div>
			}
		>
			<Input.Search
				placeholder={t('kanban.inviteSearchPlaceholder')}
				allowClear
				value={search}
				onChange={(event) => setSearch(event.target.value)}
				className="mb-4"
			/>
			<div className="max-h-[360px] overflow-y-auto">
				{filteredMembers.length === 0 ? (
					<Empty
						description={t('kanban.inviteEmpty')}
						image={Empty.PRESENTED_IMAGE_SIMPLE}
					/>
				) : (
					filteredMembers.map((item) => (
						<div
							key={item.id}
							className={`flex items-center justify-between px-3 py-2.5 border-b border-base-border cursor-pointer ${
								selected.has(item.id) ? 'bg-blue-light-1' : ''
							}`}
							onClick={() => toggleSelection(item.id)}
						>
							<div className="flex items-center gap-2.5">
								<Checkbox
									checked={selected.has(item.id)}
									onChange={() => toggleSelection(item.id)}
									onClick={(event) => event.stopPropagation()}
								/>
								<Avatar
									size={36}
									style={{ backgroundColor: 'var(--blue-medium)' }}
									className="shrink-0"
								>
									{getMemberInitial(item)}
								</Avatar>
								<div>
									<div className="font-medium">{getMemberLabel(item)}</div>
									<div className="text-xs text-base-label">{item.email}</div>
								</div>
							</div>
						</div>
					))
				)}
			</div>
		</NewModal>
	);
}
