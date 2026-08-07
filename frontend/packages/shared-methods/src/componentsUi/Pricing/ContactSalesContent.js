'use client';

import React from 'react';
import { MailOutlined, PhoneOutlined } from '@ant-design/icons';
import { displayValue } from '@repo/shared-methods/utils';
import { CONTACT_SALES_TIP, DEFAULT_SALES_EMAIL, DEFAULT_SALES_PHONE } from './constants';

/**
 * Contact Sales Popover 内容：浅蓝底与白底弹框区分；只读展示邮箱 / 电话。
 * 电话按常量原文展示（可含国家码），不走 formatPhone，避免截断。
 *
 * @param {object} [props]
 * @param {string} [props.email]
 * @param {string} [props.phone]
 * @param {string} [props.tip]
 */
const ContactSalesContent = ({
	email = DEFAULT_SALES_EMAIL,
	phone = DEFAULT_SALES_PHONE,
	tip = CONTACT_SALES_TIP,
}) => (
	<div className="w-[280px] rounded-lg bg-base-active p-3.5 shadow-md">
		<p className="m-0 text-sm font-medium leading-5 text-base-text">Contact Sales</p>
		<p className="m-0 mt-1.5 text-xs leading-5 text-base-gray">{tip}</p>
		<div className="mt-3 flex flex-col gap-2 rounded-md border border-base-border bg-white p-2.5">
			<div className="inline-flex items-center gap-2 text-sm text-base-text">
				<MailOutlined className="shrink-0 text-base-label" />
				<span className="break-all">{displayValue(email)}</span>
			</div>
			{phone ? (
				<div className="inline-flex items-center gap-2 text-sm text-base-text">
					<PhoneOutlined className="shrink-0 text-base-label" />
					{displayValue(phone)}
				</div>
			) : null}
		</div>
	</div>
);

export default ContactSalesContent;
