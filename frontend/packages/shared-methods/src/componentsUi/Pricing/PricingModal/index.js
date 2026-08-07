'use client';

import React from 'react';
import { NewModal } from '@repo/shared-methods';
import Pricing from '../index';

/**
 * 套餐定价弹框：内嵌 Pricing（variant="modal"），无 footer。
 *
 * @param {object} props
 * @param {boolean} props.open 是否展示
 * @param {() => void} props.onCancel 关闭回调
 * @param {string} [props.title] 弹框标题
 * @param {number} [props.width] 弹框宽度
 * @param {object} [props.pricingProps] 透传给 Pricing 的 props
 */
const PricingModal = ({
	open,
	onCancel,
	title = 'Adjust your plan',
	width = 960,
	pricingProps = {},
}) => {
	if (!open) {
		return null;
	}

	return (
		<NewModal title={title} width={width} onCancel={onCancel} footer={null} destroyOnHidden>
			<Pricing variant="modal" {...pricingProps} />
		</NewModal>
	);
};

export default PricingModal;
