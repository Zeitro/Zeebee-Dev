/**
 * Parse file 按钮 + 说明提示。
 * 纯展示组件：不发起请求；loading / disabled / onClick 由父组件传入。
 */

import React from 'react';
import { Button } from 'antd';
import { FileSearchOutlined } from '@ant-design/icons';
import { FieldTips } from '@repo/shared-methods';

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
 *   categoryName?: string,
 *   parseStatus?: 'idle' | 'uploaded' | 'parsing' | 'complete',
 *   onClick?: () => void,
 * }} props
 */
const ParseFileButton = ({ categoryName, parseStatus = 'idle', onClick }) => {
	const isParsing = parseStatus === 'parsing';
	const isComplete = parseStatus === 'complete';

	return (
		<div className="flex shrink-0 items-center gap-2">
			<FieldTips
				title={buildParseFileTip(categoryName)}
				isShowDescription
				rootClassName="!max-w-[360px]"
			>
				<Button
					type="primary"
					ghost
					icon={<FileSearchOutlined />}
					loading={isParsing}
					disabled={isComplete}
					onClick={onClick}
					className="mr-1"
				>
					{isParsing ? 'Parsing…' : isComplete ? 'Parsed' : 'Parse file'}
				</Button>
			</FieldTips>
		</div>
	);
};

export default ParseFileButton;
