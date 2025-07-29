import React from 'react';
import FilterField from '@/components/FilterField';
import styles from './index.module.scss';
import NewTable from '@/components/NewTable';
const LoanFilterTable = ({ columns, scroll, dataSource, filterList,rowProps }) => {
	return (
		<div className={styles.loanFilterTable}>
			<FilterField
				list={filterList}
                rowProps={rowProps}
			/>
			<div className={styles.filterAmount}></div>
			<NewTable
				dataSource={dataSource}
				columns={columns}
				scroll={scroll}
			/>
		</div>
	);
};

export default LoanFilterTable;
