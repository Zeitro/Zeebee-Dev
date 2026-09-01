/**
 * 任务描述：默认富文本预览，点击切换 WangEditor 编辑；进入编辑后保持编辑态直到组件 remount。
 */
'use client';

import dynamic from 'next/dynamic';
import { useEffect, useState } from 'react';
import { isRichTextEmpty, normalizeRichTextValue, toRichTextPreviewHtml } from '../helpers';

const EDITOR_HEIGHT_CLASS = 'h-[400px]';

/** 预加载 WangEditor chunk，减少首次进入编辑时的闪烁 */
function preloadWangEditor() {
	return import('@repo/shared-methods/components/WangEditor');
}

function EditorLoading() {
	return (
		<div
			className={`${EDITOR_HEIGHT_CLASS} overflow-hidden rounded-md border border-base-border bg-white`}
		>
			<div className="h-10 border-b border-base-border bg-base-text2 animate-pulse" />
			<div className="mx-2 mb-2 mt-2 h-[calc(100%-52px)] animate-pulse rounded bg-base-text2/60" />
		</div>
	);
}

const WangEditor = dynamic(() => preloadWangEditor(), {
	ssr: false,
	loading: EditorLoading,
});

const previewClassName = [
	'text-[15px] leading-relaxed text-base-gray break-words',
	'[&_p]:mb-2 [&_p:last-child]:mb-0',
	'[&_ul]:my-2 [&_ul]:list-disc [&_ul]:pl-5',
	'[&_ol]:my-2 [&_ol]:list-decimal [&_ol]:pl-5',
	'[&_img]:max-w-full [&_img]:rounded',
	'[&_pre]:my-2 [&_pre]:rounded [&_pre]:bg-white [&_pre]:p-2 [&_pre]:text-[13px] [&_pre]:overflow-x-auto',
	'[&_a]:text-base-primary [&_a]:underline',
].join(' ');

/**
 * @param {{
 *   value?: string,
 *   onChange?: (value: string) => void,
 *   locale?: string,
 *   placeholder?: string,
 * }} props
 */
export function TaskDescription({ value, onChange, locale, placeholder }) {
	const [isEditing, setIsEditing] = useState(false);

	useEffect(() => {
		void preloadWangEditor();
	}, []);

	const handleChange = (html) => {
		onChange?.(normalizeRichTextValue(html));
	};

	const handleStartEditing = () => {
		void preloadWangEditor();
		setIsEditing(true);
	};

	if (isEditing) {
		return (
			<WangEditor
				locale={locale}
				value={value}
				onChange={handleChange}
				wrapperProps={{ className: EDITOR_HEIGHT_CLASS }}
				editorProps={{
					defaultConfig: {
						placeholder,
					},
					style: { height: '100%', overflowY: 'hidden' },
				}}
			/>
		);
	}

	const isEmpty = isRichTextEmpty(value);

	return (
		<div
			role="button"
			tabIndex={0}
			onMouseEnter={() => void preloadWangEditor()}
			onClick={handleStartEditing}
			onKeyDown={(event) => {
				if (event.key === 'Enter' || event.key === ' ') {
					event.preventDefault();
					handleStartEditing();
				}
			}}
			className="rounded-xl border border-base-border bg-white overflow-hidden cursor-text transition-colors hover:bg-base-active/20"
		>
			<div className="min-h-[80px] px-4 py-3">
				{isEmpty ? (
					<span className="text-[13px] leading-relaxed text-base-disabled">
						{placeholder}
					</span>
				) : (
					<div
						className={previewClassName}
						dangerouslySetInnerHTML={{ __html: toRichTextPreviewHtml(value) }}
					/>
				)}
			</div>
		</div>
	);
}
