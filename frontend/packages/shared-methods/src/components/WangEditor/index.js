'use client';
import React, { useEffect, useLayoutEffect, useState, useMemo } from 'react';
import { message } from 'antd';
import { Editor, Toolbar } from '@wangeditor-next/editor-for-react';
import { Boot, i18nChangeLanguage } from '@wangeditor-next/editor';
import '@wangeditor-next/editor/dist/css/style.css';
import { uploadAwsS3 } from '../../api/common';
import { CustomVariablesMenu, CustomViewCodeMenu } from './CustomMenu';
import { MonacoModal } from './MonacoModal';

/** 业务 locale（如 next-intl 的 en/zh）→ wangeditor 语言码 */
function resolveWangEditorLocale(locale) {
	if (!locale) return 'en';
	if (locale === 'zh' || locale.startsWith('zh-')) return 'zh-CN';
	return 'en';
}

/**
 * 富文本图片走系统统一 AWS 预签名上传（uploadAwsS3），非 wangEditor 默认 multipart server。
 * @param {File} file
 * @param {(url: string, alt: string, href: string) => void} insertFn
 */
async function uploadEditorImageToAws(file, insertFn) {
	try {
		// type 固定 image：编辑器入口已是图片；后缀映射不全时仍满足 presign 字段
		const uploadedUrl = await uploadAwsS3({
			url: '/storage/presign_url',
			file,
			type: 'image',
		});
		if (!uploadedUrl) {
			message.error('Image upload failed');
			return;
		}
		insertFn(uploadedUrl, file?.name || '', uploadedUrl);
	} catch {
		message.error('Image upload failed');
	}
}

// 创建一个菜单实例管理器来存储菜单实例
let variableMenuInstance = null;
/**
 * WangEditor 富文本组件
 * @param {string} value 受控内容（HTML）
 * @param {function} onChange 内容变更回调
 * @param {string} [locale] 界面语言（如 next-intl 的 en/zh），不传默认 en
 * @param {object} editorProps 传递给 Editor 的其他属性（defaultConfig 会与内置配置 merge，勿整对象覆盖）
 * @param {object} toolbarProps 传递给 Toolbar 的其他属性
 * @param {object} wrapperProps 外层div样式/属性
 * @param {Array} variables 变量列表
 * 因wangeditor默认会调用浏览器相关方法,所以不能ssr渲染,必须用next/dynamic导入
 */
export const WangEditor = ({
	value,
	onChange,
	locale,
	editorProps = {},
	toolbarProps = {},
	wrapperProps = {},
	variables = [],
	isShowVariablesToolbar = false,
	isShowViewCodeToolbar = false,
}) => {
	const [editor, setEditor] = useState(null);
	const [showMonacoModal, setShowMonacoModal] = useState(false);
	const editorLocale = useMemo(() => resolveWangEditorLocale(locale), [locale]);
	const {
		defaultConfig: userDefaultConfig = {},
		className: editorClassName,
		...restEditorProps
	} = editorProps;
	const insertKeys = useMemo(() => {
		let keys = [];
		if (isShowVariablesToolbar) {
			keys.push({
				index: 0,
				keys: ['variable-select'], // 添加自定义菜单到工具栏
			});
		}
		if (isShowViewCodeToolbar) {
			keys.push({
				index: 0,
				keys: ['monaco-code-modal-view'], // 添加自定义菜单到工具栏
			});
		}
		return keys;
	}, [isShowVariablesToolbar, isShowViewCodeToolbar]);
	// 默认编辑器配置（与 userDefaultConfig merge；MENU_CONF.uploadImage 默认走 AWS）
	const defaultEditorConfig = {
		placeholder: 'Please enter content...',
		language: editorLocale,
		styleWithCSS: true,
		...userDefaultConfig,
		MENU_CONF: {
			...userDefaultConfig.MENU_CONF,
			uploadImage: {
				async customUpload(file, insertFn) {
					await uploadEditorImageToAws(file, insertFn);
				},
				...userDefaultConfig.MENU_CONF?.uploadImage,
			},
		},
	};

	// 默认工具栏配置
	const defaultToolbarConfig = {
		...toolbarProps,
		excludeKeys: [...(toolbarProps.excludeKeys || []), 'group-video'],
	};

	// 变更回调
	const handleChange = (editor) => {
		const newHtml = editor.getHtml();
		onChange?.(newHtml);
	};

	const onCreated = (editor) => {
		setEditor(editor);
		let menus = [];
		console.log(isShowVariablesToolbar);

		if (isShowVariablesToolbar) {
			menus.push({
				key: 'variable-select',
				factory() {
					variableMenuInstance = new CustomVariablesMenu();
					return variableMenuInstance;
				},
			});
		}
		if (isShowViewCodeToolbar) {
			menus.push({
				key: 'monaco-code-modal-view',
				factory() {
					return new CustomViewCodeMenu(() => setShowMonacoModal(true));
				},
			});
		}

		if (!menus?.length) return;
		if (menus?.length > 1) {
			Boot.registerModule({
				menus,
			});
		} else {
			Boot.registerMenu(menus[0]);
		}
	};
	//     useEffect(() => {
	//   let menus = [];
	//   if (isShowVariablesToolbar) {
	//     menus.push({
	//       key: 'variable-select',
	//       factory() {
	//         variableMenuInstance = new CustomVariablesMenu();
	//         return variableMenuInstance;
	//       },
	//     });
	//   }
	//   if (isShowViewCodeToolbar) {
	//     menus.push({
	//       key: 'monaco-modal-btn',
	//       factory() {
	//         return new CustomViewCodeMenu(() => setShowMonacoModal(true));
	//       },
	//     });
	//   }
	//   if (menus.length > 1) {
	//     Boot.registerModule({ menus });
	//   } else if (menus.length === 1) {
	//     Boot.registerMenu(menus[0]);
	//   }
	// }, [isShowVariablesToolbar, isShowViewCodeToolbar]);

	useLayoutEffect(() => {
		i18nChangeLanguage(editorLocale);
	}, [editorLocale]);

	useEffect(() => {
		return () => {
			editor?.destroy?.();
		};
	}, [editor]);

	// 当 variables 更新时，直接更新菜单选项
	useEffect(() => {
		if (variableMenuInstance && isShowVariablesToolbar) {
			// 直接更新菜单实例的变量数组
			variableMenuInstance.updateVariables(variables);
		}
	}, [variables, isShowVariablesToolbar]);

	return (
		<>
			<div {...wrapperProps} className={['flex flex-col', wrapperProps?.className].join(' ')}>
				<Toolbar
					editor={editor}
					defaultConfig={{
						...defaultToolbarConfig,
						insertKeys,
					}}
					mode="default"
					className={[
						'border border-base-border border-b-0 rounded-t-md flex-shrink-0',
						toolbarProps?.className,
					].join(' ')}
					excludeKeys={['group-video']}
					{...toolbarProps}
				/>
				<Editor
					key={editorLocale}
					value={value}
					onCreated={onCreated}
					onChange={handleChange}
					defaultConfig={defaultEditorConfig}
					mode="default"
					className={[
						'flex-1 overflow-hidden min-h-[160px] border border-base-border rounded-b-md',
						editorClassName,
					].join(' ')}
					{...restEditorProps}
				/>
			</div>
			{showMonacoModal && <MonacoModal setOpen={setShowMonacoModal} data={value} />}
		</>
	);
};

export default WangEditor;
