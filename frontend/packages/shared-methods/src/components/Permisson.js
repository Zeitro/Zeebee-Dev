// packages/shared-methods/src/components/Permission.js
import React, { isValidElement, cloneElement } from 'react';
import useAuth from '../hooks/useAuth';
import { message } from 'antd';

/**
 * 权限指令组件
 * @param {string|string[]} code - 权限标识或权限标识数组
 * @param {string} type - 无权限时的处理方式：'disabled' | 'hidden' | 'none'
 * @param {React.ReactNode} children - 子组件
 * @param {React.ReactNode} fallback - 无权限时的替代内容
 */
const Permission = ({ code, type = 'hidden', children, fallback = null, ...restProps }) => {
	const { hasPermission, hasAnyPermission } = useAuth();

	// 判断是否有权限
	const hasAuth = Array.isArray(code) ? hasAnyPermission(code) : hasPermission(code);

	// 有权限，直接渲染子组件
	if (hasAuth) {
		return children;
	}

	// 无权限时隐藏
	if (type === 'hidden') {
		return fallback;
	}

	// 无权限时禁用（克隆子组件并添加 disabled 属性）
	if (type === 'disabled') {
		return React.Children.map(children, (child) => {
			if (React.isValidElement(child)) {
				return React.cloneElement(child, {
					disabled: true,
					...restProps,
				});
			}
			return child;
		});
	}
	if (type === 'quote') {
		const handleClick = (e) => {
			e.preventDefault();
			e.stopPropagation();
			message.open({
				type: 'warning',
				content: 'Please pay to activate this permission before you can use it.',
			});
		};
		// 处理 React 元素，包装点击事件
		if (isValidElement(children)) {
			return cloneElement(children, {
				onClick: (e) => {
					handleClick(e);
					// 保留子组件原有 onClick（但不会执行，因为已阻断）
					if (children.props.onClick) {
						children.props.onClick(e);
					}
				},
				style: {
					...children.props.style,
					cursor: 'not-allowed',
					opacity: 0.7, // 视觉提示
				},
			});
		}
		// 非 React 元素，用 div 包装
		return (
			<div
				onClick={handleClick}
				style={{
					display: 'inline-block',
					cursor: 'not-allowed',
					opacity: 0.7,
				}}
				title="No permission"
			>
				{children}
			</div>
		);
	}

	// 无权限时渲染替代内容
	return fallback;
};

export default Permission;
