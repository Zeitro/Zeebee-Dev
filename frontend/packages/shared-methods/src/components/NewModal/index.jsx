import { Modal, Button } from 'antd';
import React from 'react';
import styles from './index.module.scss';
export const NewModal = ({ children, customFooter = false, confirmLoading = false, ...props }) => {
	return (
		<Modal
			open={true}
			centered
			maskClosable={false}
			confirmLoading={confirmLoading}
			okText="Save"
			{...props}
			okButtonProps={{
				size: 'middle',
				...(props?.okButtonProps || {}),
			}}
			cancelButtonProps={{
				size: 'middle',
				...(props?.cancelButtonProps || {}),
			}}
			className={[
				styles.NewModal,
				// !customFooter && styles.NewModal_padding,
				props?.className,
			].join(' ')}
		>
			{children}
			{!props?.footer && customFooter && (
				<div className="flex items-center justify-between p-y-4 y-x-7 border-t border-base-border">
					{customFooter?.cancel && (
						<Button type="link" onClick={props?.onCancel}>
							{props?.cancelText || 'Cancel'}
						</Button>
					)}
					{customFooter?.save && (
						<Button type="primary" onClick={props?.onOk}>
							{props?.okText || 'Save'}
						</Button>
					)}
				</div>
			)}
		</Modal>
	);
};
