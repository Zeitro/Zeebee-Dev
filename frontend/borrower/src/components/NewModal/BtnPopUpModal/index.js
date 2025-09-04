'use client';
import { Button, Modal } from 'antd';
import { useState } from 'react';
import { ExclamationCircleOutlined } from '@ant-design/icons';
import { title } from 'process';
/**
 *	按钮 + 确认弹窗
 * @param {*} type confirm | warning | info | error
 */
export const BtnPopUpModal = ({
	children,
	modalType = '',
	config = {},
	getContent = () => {}, // 弹框组件
	...props
}) => {
	const [modal, contextHolder] = Modal.useModal();
	const [open, setOpen] = useState('');
	const handleClick = () => {
		if (modalType) {
			modal[modalType]?.({
				title: <div className="font-medium">{title}</div>,
				icon: <ExclamationCircleOutlined />,
				content: null,
				maskClosable: false,
				centered: true,
				...config,
				okButtonProps: {
					size: 'middle',
					className: 'border-[6px]',
					...(config?.okButtonProps || {}),
				},
				cancelButtonProps: {
					size: 'middle',
					className: 'border-[6px]',
					...(config?.cancelButtonProps || {}),
				},
			});
			return;
		}
		setOpen(true);
	};

	return (
		<>
			<Button type="link" size="small" onClick={handleClick} {...props}>
				{children}
			</Button>
			{contextHolder}
			{open && getContent(setOpen)}
		</>
	);
};

// 删除组件
export const DeleteBtn = ({ title = 'Are you sure you want to delete this?', ...props }) => {
	return (
		<BtnPopUpModal
			danger={true}
			modalType="confirm"
			config={{ title, okText: 'Delete', okButtonProps: { danger: true }, ...props }}
		>
			Delete
		</BtnPopUpModal>
	);
};
