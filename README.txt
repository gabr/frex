frex: File REgeX replacement

Usage:
  frex regex_pattern_to_replace value_to_replace_to file_path1 file_path2...

Example:

> frex .*ReplaceMe Line2 file1.txt

file1.txt before change    >    file1.txt after change
  Line1                    >      Line1
  LineReplaceMe            >      Line2
  Line3                    >      Line3
