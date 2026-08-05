'use client';

/**
 * Parse file 按钮：说明提示 + 触发解析请求。
 * POST .../parse → request_id，写入 transaction store，再调用 DocumentsTab 注入的 pollParseStatus。
 */

import React, { useMemo, useState } from 'react';
import { Button, message } from 'antd';
import { FileSearchOutlined } from '@ant-design/icons';
import { throttle } from 'lodash';
import { FieldTips } from '@repo/shared-methods';
import { createRequest } from '@repo/shared-methods/utils';

/** 与 GET .../parse/result 返回的 status 对齐；Finished 为成功终态 */
const PARSE_SUCCESS_STATUS = 'Finished';

/**
 * @param {string} [categoryName]
 * @returns {React.ReactNode}
 */
function buildParseFileTip(categoryName) {
	return (
		<div className="flex flex-col gap-1">
			<span className="text-base">
				Parses all files under {categoryName ?? 'this category'}.
			</span>
			<span>
				Related files are linked together — please upload all related files before clicking
				Parse file.
			</span>
			<span>
				Please first upload the files with important information to extract, such as
				Purchase price.
			</span>
		</div>
	);
}

/**
 * @param {{
 *   category?: object,
 *   documentId?: number,
 *   transactionId?: number|string,
 *   pollParseStatus?: (transactionId: number|string, documentId: number|string) => Promise<object|null>,
 * }} props
 */
const ParseFileButton = ({ documentId, transactionId, category, pollParseStatus }) => {
	/** @type {['idle' | 'uploaded' | 'parsing' | 'complete', Function]} */
	const [parseStatus, setParseStatus] = useState('idle');
	const { request } = useMemo(() => createRequest(), []);

	const isParsing = parseStatus === 'parsing';
	const parser = category?.file_parse_config?.rules?.find(
		(rule) => rule.triggered_by === 'Manual',
	)?.parser;

	/**
	 * POST /borrower_property_transaction/{id}/parse → request_id。
	 * 节流挡连点；parsing 时按钮 loading，双重避免重复触发。
	 */
	const handleParseFile = throttle(
		async () => {
			try {
				setParseStatus('parsing');
				const params = {
					document_id: documentId,
					parser,
				};
				const res = await request(`/borrower_property_transaction/${transactionId}/parse`, {
					method: 'POST',
					body: params,
				});
				params.request_id = res?.request_id;
				if (!params.request_id) {
					message.error('Parse started but no request_id returned');
					setParseStatus('idle');
					return;
				}

				const result = await pollParseStatus?.(transactionId, documentId, params);
				if (result?.status === PARSE_SUCCESS_STATUS) {
					setParseStatus('complete');
					return;
				}
				setParseStatus('idle');
			} catch {
				setParseStatus('idle');
			}
		},
		1000,
		{ trailing: false },
	);

	return (
		<div className="flex shrink-0 items-center gap-2">
			<FieldTips
				title={buildParseFileTip(category?.name)}
				isShowDescription
				rootClassName="!max-w-[360px]"
			>
				<Button
					type="primary"
					ghost
					icon={<FileSearchOutlined />}
					loading={isParsing}
					onClick={handleParseFile}
					className="mr-1"
				>
					Parse all files
				</Button>
			</FieldTips>
		</div>
	);
};

export default ParseFileButton;
