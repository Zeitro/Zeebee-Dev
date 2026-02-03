'use client';
import React, { useRef, useEffect, useState, useMemo } from 'react';
import { Editor, Toolbar } from '@wangeditor-next/editor-for-react';
import { Boot, i18nChangeLanguage } from '@wangeditor-next/editor';
import '@wangeditor-next/editor/dist/css/style.css';
import { CustomVariablesMenu, CustomViewCodeMenu } from './CustomMenu';
import { MonacoModal } from './MonacoModal';
i18nChangeLanguage('en');
// 创建一个菜单实例管理器来存储菜单实例
let variableMenuInstance = null;
/**
 * WangEditor 富文本组件
 * @param {string} value 受控内容（HTML）
 * @param {function} onChange 内容变更回调
 * @param {object} editorProps 传递给 Editor 的其他属性
 * @param {object} toolbarProps 传递给 Toolbar 的其他属性
 * @param {object} wrapperProps 外层div样式/属性
 * @param {Array} variables 变量列表
 * 因wangeditor默认会调用浏览器相关方法,所以不能ssr渲染,必须用next/dynamic导入
 */
export const WangEditor = ({
	value,
	onChange,
	editorProps = {},
	toolbarProps = {},
	wrapperProps = {},
	variables = [],
	isShowVariablesToolbar = false,
	isShowViewCodeToolbar = false,
}) => {
	const [editor, setEditor] = useState(null);
	const initRef = useRef(true);
	const [showMonacoModal, setShowMonacoModal] = useState(false);
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
	// 默认编辑器配置
	const defaultEditorConfig = {
		placeholder: 'Please enter content...',
		language: 'en', // 设置为英语
		styleWithCSS: true,
		// 其他配置
		...editorProps.defaultConfig,
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

	// 在编辑器创建后注册自定义菜单
	useEffect(() => {
		// if (editor && initRef.current && isShowVariablesToolbar) {
		// 	// 创建菜单类
		// 	const variableConf = {
		// 		key: 'variable-select',
		// 		factory() {
		// 			variableMenuInstance = new CustomVariablesMenu();
		// 			return variableMenuInstance;
		// 		},
		// 	};

		// 	Boot.registerMenu(variableConf);
		// 	initRef.current = false;
		// }
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
					value={value}
					onCreated={onCreated}
					onChange={handleChange}
					defaultConfig={defaultEditorConfig}
					mode="default"
					className={[
						'flex-1 overflow-hidden min-h-[160px] border border-base-border rounded-b-md',
						editorProps?.className,
					].join(' ')}
					{...editorProps}
				/>
			</div>
			{showMonacoModal && <MonacoModal setOpen={setShowMonacoModal} data={value} />}
		</>
	);
};

export default WangEditor;
