'use client';
import React, { useMemo, useRef, useState } from 'react';
import { createRequest, formatPhone } from '@repo/shared-methods/utils';
import { LinkOutlined, ReloadOutlined } from '@ant-design/icons';
import { FilterTableCard } from '@repo/shared-methods';
// import useContent from './components/useContent';
// import MainModal from './components/Main';
import PersonSvg from '@repo/shared-methods/assets/images/person.svg';
import EmailSvg from '@repo/shared-methods/assets/images/email.svg';
import PhoneSvg from '@repo/shared-methods/assets/images/phone.svg';
import { Image } from 'antd';
import Link from 'next/link';

const CONTACT_LENDER_API = '/organization/lender/page';
const ContactLender = () => {
	const { cancelMatch } = useMemo(() => createRequest(), []);
	const tableFilterRef = useRef();
	const [searchValue, setSearchValue] = useState('');
	// const { columns, handleDelete, loading } = useContent({
	// 	handleSearch: () => tableFilterRef?.current?.handleSearch(),
	// });

	const filterList = useMemo(
		() => [
			{
				filterType: 'input',
				placeholder: 'Search by any field',
				paramsKey: 'keyword',
				value: searchValue,
				onChange: ({ target: { value } }) => setSearchValue(value),
			},
			{
				filterType: 'button',
				onClick: () => {
					cancelMatch(CONTACT_LENDER_API);
					tableFilterRef?.current?.handleSearch();
				},
			},
			{
				filterType: 'button',
				type: 'default',
				text: 'Reset',
				icon: <ReloadOutlined />,
				onClick: () => {
					cancelMatch(CONTACT_LENDER_API);
					setSearchValue('');
					tableFilterRef?.current?.handleSearch({
						filter: {
							keyword: '',
						},
					});
				},
			},
		],
		[searchValue, cancelMatch],
	);
	return (
		<FilterTableCard
			filterWrapClassName="px-8 pt-10 pb-8 flex-1 overflow-hidden"
			ref={tableFilterRef}
			filterList={filterList}
			hasFilterField={true}
			// pagination={pagination}
			api={CONTACT_LENDER_API}
			typeKey="card"
			itemClassName="w-[300px]"
			cardTitle="Contact lender"
			renderItem={(row) => {
				return (
					<>
						<div className="flex items-center">
							<Image
								key={row?.id}
								src={row?.logo || null}
								width={48}
								height={48}
								alt="logo"
								preview={false}
								fallback={'/images/noOrgLogo.svg'}
								style={{ objectFit: 'contain', borderRadius: 8 }}
							/>

							<p className="font-medium text-base text-black-dark pl-2 overflow-hidden">
								{row?.name || '--'}
							</p>

							{/* <div className="flex items-center">
								<BtnPopUpModal
									getContent={(setOpen) => (
										<MainModal
											setOpen={setOpen}
											handleSearch={() =>
												tableFilterRef?.current?.handleSearch()
											}
											rowData={row}
										/>
									)}
									type="link"
									size="large"
									// isButton={false}
									// icon={<EditOutlined />}
									className="!px-2"
								>
									Edit
								</BtnPopUpModal>
								<DeleteBtn
									type="link"
									title="Are you sure delete this borrower?"
									content={row.full_name}
									// isButton={false}
									onOk={() => handleDelete(row.id)}
									loading={loading}
									danger
									className="!px-2"
								>
									Delete
								</DeleteBtn>
							</div> */}
						</div>

						<div className="flex items-center gap-4 pt-3 overflow-hidden">
							<PersonSvg className="flex-shrink-0" />
							<span className="text-sm/6 common-ellipsis">
								{row?.contact_name || '--'}
							</span>
						</div>
						<div className="flex items-center gap-4 pt-3 overflow-hidden">
							<EmailSvg className="flex-shrink-0" />
							<span className="text-sm/6 common-ellipsis">
								{row?.contact_email || '--'}
							</span>
						</div>
						<div className="flex items-center gap-4 pt-3">
							<PhoneSvg className="flex-shrink-0" />
							<span className="text-sm/6">
								{formatPhone(row?.contact_phone) || '--'}
							</span>
						</div>
						<div className="flex items-center gap-6 pt-3">
							<LinkOutlined
								className="flex-shrink-0"
								style={{
									fontSize: 18,
									color: 'var(--base-primary)',
									marginLeft: 2,
								}}
							/>
							{row?.website ? (
								<Link
									href={row?.website || ''}
									target="_blank"
									className="overflow-hidden"
								>
									<span className="text-sm/6 common-ellipsis block ">
										{row?.website}
									</span>
								</Link>
							) : (
								<span className="text-sm/6">--</span>
							)}
						</div>
					</>
				);
			}}
		/>
	);
};
export default ContactLender;
