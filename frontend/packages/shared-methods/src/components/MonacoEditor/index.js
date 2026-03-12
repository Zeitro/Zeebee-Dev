import React, { useEffect, useState } from 'react';
import MonacoEditor, { useMonaco } from '@monaco-editor/react';
import prettier from 'prettier/standalone';
import parserHtml from 'prettier/plugins/html'; // 新版写法

// 防止重复加载 monaco-css
let cssLoaded = false;

/**
 * 通用 Monaco 代码编辑器组件，支持 CSS/HTML/JS 代码提示、html自动格式化、行内css智能提示
 */
const CodeEditor = ({
	value = '',
	onChange,
	language = 'javascript',
	options = {},
	editorProps = {},
	height = '60vh',
	...rest
}) => {
	const monaco = useMonaco();
	const [editor, setEditor] = useState(null);
	const [hasFormatted, setHasFormatted] = useState(false);

	useEffect(() => {
		if (monaco && !cssLoaded) {
			cssLoaded = true;
			import('monaco-editor/esm/vs/language/css/monaco.contribution');
		}
		if (monaco) {
			// html 自动格式化，使用 prettier
			monaco.languages.registerDocumentFormattingEditProvider('html', {
				provideDocumentFormattingEdits: (model) => {
					let code = model.getValue();
					// 预处理：去除自闭合标签的闭合部分，避免prettier抛错
					// 例如 <col></col>、<br></br>、<hr></hr>、<img></img> 等
					code = code.replace(
						/<(col|br|hr|img|input|meta|link|source|track|wbr|area|base|embed|param|command|keygen|menuitem)\b[^>]*><\/\1>/gi,
						(m, tag) => {
							return m.replace(new RegExp(`<\\/${tag}>`, 'i'), '');
						},
					);
					try {
						const formatted = prettier.format(code, {
							parser: 'html',
							plugins: [parserHtml],
							tabWidth: 2,
							useTabs: false,
							printWidth: 120,
						});
						return [
							{
								range: model.getFullModelRange(),
								text: formatted,
							},
						];
					} catch {
						return [];
					}
				},
			});
			// html 文件内 style="..." 时，启用 monaco-css 的补全
			monaco.languages.registerCompletionItemProvider('html', {
				triggerCharacters: ['"', ';', ':', ' '],
				provideCompletionItems: (model, position) => {
					const textUntilPosition = model.getValueInRange({
						startLineNumber: 1,
						startColumn: 1,
						endLineNumber: position.lineNumber,
						endColumn: position.column,
					});
					// 检查是否在 style="..." 内
					const styleAttrMatch = /style\s*=\s*["'][^"'>]*$/i.exec(textUntilPosition);
					if (styleAttrMatch && monaco.languages.css) {
						// 触发 monaco-css 的补全（仅属性名，值等更复杂可扩展）
						const cssProps =
							monaco.languages.css.cssDefaults?._options?._data?.properties || [];
						return {
							suggestions: cssProps.map((prop) => ({
								label: prop.name,
								kind: monaco.languages.CompletionItemKind.Property,
								insertText: prop.name + ': ',
								detail: prop.description,
							})),
						};
					}
					return { suggestions: [] };
				},
			});
		}
	}, [monaco]);

	// 内容ready后自动格式化，仅触发一次
	// useEffect(() => {
	// 	if (
	// 		editor &&
	// 		language === 'html' &&
	// 		!hasFormatted &&
	// 		value &&
	// 		typeof value === 'string' &&
	// 		value.trim().length > 0
	// 	) {
	// 		setHasFormatted(true);
	// 		setTimeout(() => {
	// 			editor
	// 				.getAction('editor.action.formatDocument')
	// 				.run()
	// 				.catch(() => {});
	// 		}, 1000); // 延迟触发，确保内容已渲染
	// 	}
	// }, [language, hasFormatted, editor, value]);

	const handleEditorDidMount = (editor, monacoInstance) => {
		setEditor(editor);
	};

	return (
		<MonacoEditor
			value={value}
			language={language}
			onChange={onChange}
			options={{
				fontSize: 14,
				minimap: { enabled: false },
				scrollBeyondLastLine: false,
				wordWrap: 'on',
				automaticLayout: true,
				...options,
			}}
			height={height}
			onMount={handleEditorDidMount}
			theme="vs-dark"
			{...editorProps}
			{...rest}
		/>
	);
};

export default CodeEditor;
