# the sample sec file for fake wechat
%define _topdir /home/aaa/rpmbuild/
%define name fakewechat
%define release 1
%define version 1.1
%define buildroot %{_topdir}/BUILDROOT/%{name}-%{version}


BuildRoot: %{buildroot}
Name: %{name}		
Version:  %{version}	
Release:  %{release}
Summary: A fake wechat Demo	

Group:	 Development/Tools	
License: COPY RIGHT	
URL:  http://www.gitgub.com/xiaojiaqi/		
Source0: %{name}-%{version}.tar.gz	

Prefix: /usr/local/fakewechat

#BuildRequires:	
#Requires:	

%description
A simple wechat demo

%prep
%setup -q
echo "I will start"

%build
#%configure
#make %{?_smp_mflags}
./configure
g++ -o hello hello.cpp

%install
mkdir -p %{buildroot}/usr/local/%{name}
#cp hello %{buildroot}
cp hello %{buildroot}/usr/local/%{name}

%files
#%doc
%defattr (-,root,root)

/usr/local/fakewechat/hello
#/hello
#/usr/bin/${name}

%changelog

