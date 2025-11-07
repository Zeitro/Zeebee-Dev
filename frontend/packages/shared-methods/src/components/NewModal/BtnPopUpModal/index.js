'use client';
import { Button, Modal } from 'antd';
import { useState } from 'react';
import { ExclamationCircleOutlined } from '@ant-design/icons';
// 弹框的初始配置样式
export const initConfigFn = (config = {}) => {
	return {
		icon: <ExclamationCircleOutlined />,
		content: null,
		maskClosable: false,
		centered: true,
		confirmLoading: false,
		...config,
		title: <div className="font-medium">{config.title}</div>,
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
	};
};
/**
 *	按钮 + 确认弹窗
 * @param {*} type confirm | warning | info | error
 */
export const BtnPopUpModal = ({
	children,
	modalType = '',
	config = {},
	getContent = () => {}, // 弹框组件
	isButton = true,
	predicate = () => true, // 用于条件判断不通过,则不弹出弹框
	...props
}) => {
	const [modal, contextHolder] = Modal.useModal();
	const [open, setOpen] = useState('');
	const handleClick = () => {
		const isPass = predicate?.();
		if (!isPass) return;
		if (modalType) {
			const initValue = initConfigFn(config);
			modal[modalType]?.(initValue);
			return;
		}
		setOpen(true);
	};

	return (
		<>
			{isButton ? (
				<Button type="link" size="small" onClick={handleClick} {...props}>
					{children}
				</Button>
			) : (
				<div onClick={handleClick} className={props?.className}>
					{children}
				</div>
			)}
			{contextHolder}
			{open && getContent(setOpen)}
		</>
	);
};

// 删除组件
export const DeleteBtn = ({
	title = 'Are you sure you want to delete this?',
	children,
	isButton,
	icon,
	type,
	danger = true,
	size,
	disabled,
	className,
	predicate = () => true,
	...props
}) => {
	return (
		<BtnPopUpModal
			danger={danger}
			modalType="confirm"
			config={{ title, okText: 'Delete', okButtonProps: { danger: true }, ...props }}
			isButton={isButton}
			icon={icon}
			type={type}
			size={size}
			disabled={disabled}
			className={className}
			predicate={predicate}
		>
			{children || 'Delete'}
		</BtnPopUpModal>
	);
};

// 离开并保存草稿
export const LeaveAndSaveDraftBtn = ({
	handleSaveDraft,
	isChange,
	disabled,
	router,
	confirmLoading,
}) => {
	return (
		<BtnPopUpModal
			config={{
				onOk: () => handleSaveDraft(),
				title: 'Tips',
				content: `${isChange ? 'The data is not saved yet, ' : ''}are you sure you want to leave the current page?`,
				okText: 'Leave and save draft',
				cancelText: 'Leave',
				closable: true,
				cancelButtonProps: {
					onClick: () => router.back(),
				},
				okButtonProps: {
					disabled: !isChange,
				},
				confirmLoading,
			}}
			modalType="confirm"
			type="default"
			size="middle"
			disabled={disabled}
		>
			Cancel
		</BtnPopUpModal>
	);
};
