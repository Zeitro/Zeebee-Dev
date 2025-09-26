'use client';
import React from 'react';
import ApplicantInfo from '../../components/CreateZeeBee/ApplicantInfo';
import BorrowerInfo from '../../components/CreateZeeBee/BorrowerInfo';
import PropertyInfo from '../../components/CreateZeeBee/PropertyInfo';
import Loan from '../../components/CreateZeeBee/Loan';
import Experience from '../../components/CreateZeeBee/Experience';
import ThirdPartyService from '../../components/CreateZeeBee/ThirdPartyService';
import Declaration from '../../components/CreateZeeBee/Declaration';
import Demographics from '../../components/CreateZeeBee/Demographics';
import ReviewAndSign from '../../components/CreateZeeBee/ReviewAndSign';

const CommonZeeBeeForm = ({ form, pageKey, children, setIsChange, rowData }) => {
	const getContent = () => {
		switch (pageKey) {
			case 'borrowerInformation':
				return (
					<>
						<ApplicantInfo form={form} />
						<BorrowerInfo form={form} />
						{children}
					</>
				);
			case 'propertyAndLoan':
				return (
					<>
						<PropertyInfo form={form} />
						<Loan form={form} />
						{children}
					</>
				);
			case 'experienceAndThirdService':
				return (
					<>
						<Experience form={form} rowData={rowData} setIsChange={setIsChange} />
						<ThirdPartyService
							form={form}
							rowData={rowData}
							setIsChange={setIsChange}
						/>
						{children}
					</>
				);
			case 'declarationAndDemographics':
				return (
					<>
						<Declaration form={form} />
						<Demographics form={form} />
						{children}
					</>
				);
			case 'reviewAndSign':
				return (
					<>
						<ReviewAndSign form={form} />
						{children}
					</>
				);
			default:
				return '';
		}
	};
	return getContent();
};
export default CommonZeeBeeForm;
