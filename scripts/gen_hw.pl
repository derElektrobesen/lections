#!/usr/bin/perl

use strict;
use warnings;

use Cwd qw( getcwd );
use POSIX qw( strftime );

srand time;

my $students_list = shift
	or die "Usage: $0 students_list.txt\n";

open my $f, '<', $students_list
	or die "Can't open students list $students_list: $!\n";

my $vars_dir = getcwd . "/variants";

my $workdir = getcwd . "/result/" . strftime("%F.%H.%M.%S", localtime);
mkdir $workdir
	or die "Can't create dir $workdir: $!\n";

my $logf = "$workdir/$students_list.log";
open my $log, '>', $logf
	or die "Can't open $logf: $!\n";

open my $templatef, '<', "$vars_dir/main.tex"
	or die "Can't open main.tex into $vars_dir: $!\n";

my $template;
{
	local $/;
	$template = <$templatef>;
}
close $templatef;

while (<$f>) {
	my $tmpl = $template;
	my ($surname, $name) = split /\s+/;

	print $log "$name $surname";

	$tmpl =~ s/%STUDENT%/$name $surname/
		or die "STUDENT macro not found in template\n";

	my %used_uniqs;
	while (my ($match, $varname) = $tmpl =~ /(%VAR:([^%]+)%)/) {
		my $var = gen_var($varname);
		next if @used_uniqs{@{$var->{uniqs}}}; # retry

		$used_uniqs{$_} = 1
			for @{$var->{uniqs}};

		$tmpl =~ s/$match/$var->{path}/;

		print $log "\t$var->{var}";
	}

	print $log "\n";

	my $outfname = "$workdir/main";
	open my $out, '>', "$outfname.tex"
		or die "Can't write to $outfname: $!\n";

	print $out $tmpl;

	close $out;

	my $pdfname = "$workdir/$surname\_$name.pdf";
	`cd $workdir && xelatex $outfname.tex`
		or die "Can't create $pdfname. Run `xelatex $outfname.tex` for more info\n";

	rename "$outfname.pdf", $pdfname
		or die "Can't rename $outfname.pdf: $!\n";
}

close $log;
close $f;

my %dirs_cache;
sub gen_var {
	my $var_dir = shift;
	unless ($dirs_cache{$var_dir}) {
		my $d = lc "$vars_dir/$var_dir";
		opendir my $dh, $d
			or die "Can't open $d: $!\n";

		my @files = map {
				parse_var($d, $_)
			} grep {
				/\.tex$/ && -f "$d/$_"
			} readdir $dh;

		die "Empty directory: $d\n"
			unless @files;

		closedir $dh;

		$dirs_cache{$var_dir} = \@files;
	}

	return $dirs_cache{$var_dir}[rand @{$dirs_cache{$var_dir}}];
}

sub parse_var {
	my ($dirpath, $fname) = @_;
	my $filename = "$dirpath/$fname";

	open my $f, '<', $filename
		or die "Can't open var file $filename: $!\n";

	local $/;
	my $tmpl = <$f>;
	close $f;

	my @uniqs;
	for ($tmpl =~ /%UNIQ:([^%]+)%/g) {
		push @uniqs, lc;
	}

	return {
		path => $filename,
		var => $fname,
		uniqs => \@uniqs,
	};
}
