/**
 * 多执行者芯片：Form.Item 受控 value/onChange；选项来自项目成员（不含 viewer）。
 */
'use client';

import { CheckOutlined, CloseOutlined, PlusOutlined, UserOutlined } from '@ant-design/icons';
import { Avatar, Dropdown } from 'antd';
import { displayValue } from '@repo/shared-methods/utils';
import { getMemberInitial, getMemberLabel } from './helpers';

/**
 * @param {{
 *   value?: number[],
 *   onChange?: (ids: number[]) => void,
 *   options?: Array,
 *   fallbacks?: Array,
 * }} props
 */
export function AssigneeChips({ value, onChange, options = [], fallbacks = [] }) {
	const ids = value || [];
	const findMember = (id) =>
		options.find((member) => member.id === id) || fallbacks.find((member) => member.id === id);

	const menuItems = options.map((member) => {
		const name = getMemberLabel(member);
		const selected = ids.includes(member.id);
		return {
			key: String(member.id),
			label: (
				<span className="inline-flex items-center gap-2 min-w-[140px]">
					<Avatar
						size={20}
						style={{ backgroundColor: 'var(--blue-medium)', fontSize: 10 }}
					>
						{getMemberInitial(member)}
					</Avatar>
					{displayValue(name)}
					{selected ? <CheckOutlined className="text-base-primary ml-auto" /> : null}
				</span>
			),
			onClick: () =>
				onChange?.(selected ? ids.filter((id) => id !== member.id) : [...ids, member.id]),
		};
	});

	return (
		<div className="flex items-center gap-1.5 flex-wrap">
			{ids.map((id) => {
				const member = findMember(id);
				const name = getMemberLabel(member, id);
				return (
					<span key={id} className="relative inline-flex">
						<span className="inline-flex items-center gap-1.5 bg-base-text2 border border-base-border rounded-2xl py-0.5 pr-2.5 pl-1">
							<Avatar
								size={22}
								style={{ backgroundColor: 'var(--blue-medium)', fontSize: 11 }}
							>
								{getMemberInitial(member, id)}
							</Avatar>
							<span className="font-medium text-[13px]">{displayValue(name)}</span>
						</span>
						<button
							type="button"
							className="absolute -top-1.5 -right-1.5 w-4 h-4 rounded-full bg-white border border-base-border inline-flex items-center justify-center z-[1]"
							onClick={(event) => {
								event.stopPropagation();
								onChange?.(ids.filter((item) => item !== id));
							}}
						>
							<CloseOutlined className="text-[8px] text-base-label" />
						</button>
					</span>
				);
			})}
			<Dropdown menu={{ items: menuItems }} trigger={['click']} placement="bottomLeft">
				<button
					type="button"
					className="w-[26px] h-[26px] rounded-full bg-base-text2 inline-flex items-center justify-center shrink-0"
				>
					{menuItems.length ? (
						<PlusOutlined className="text-xs text-base-gray" />
					) : (
						<UserOutlined className="text-xs text-base-gray" />
					)}
				</button>
			</Dropdown>
		</div>
	);
}
