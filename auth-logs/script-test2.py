import re
from collections import defaultdict
from datetime import datetime
import os
import argparse
from typing import Dict, List, Tuple
import time
from plyer import notification
class SSHLogAnalyzer:
    def __init__(self, log_path: str):
        self.log_path = log_path
        self.attempts = []
        self.user_attempts=defaultdict(list)
    def parse_log_file(self) -> None:
        """Parse the SSH log file and store attempts."""
        try:
            with open(self.log_path, 'r') as file:
                for line in file:
                    if 'sshd' in line and ('Failed password' in line or 
                                         'Invalid user' in line or 
                                         'Accepted password' in line or 
                                         'Accepted publickey' in line):
                        self._process_line(line)
        except FileNotFoundError:
            print(f"Error: Log file '{self.log_path}' not found.")
            exit(1)

    def _process_line(self, line: str) -> None:
        """Process each line and extract login attempt information."""
        timestamp = self._extract_timestamp(line)
        ip = self._extract_ip(line)
        username = self._extract_username(line)

        if not all([timestamp, ip, username]):
            return

        # Detailed status detection
        if "Failed password" in line:
            if "invalid user" in line.lower():
                status = "Failed Login (Invalid User)"
            else:
                status = "Failed Login (Wrong Password)"
        elif "Invalid user" in line:
            status = "Failed Login (Invalid User)"
        elif "Accepted password" in line:
            status = "Successful Login (Password)"
        elif "Accepted publickey" in line:
            status = "Successful Login (Public Key)"
        else:
            status = "Failed Login"

        # Convert timestamp to UTC if possible
        try:
            dt = datetime.strptime(timestamp, "%b %d %H:%M:%S")
            dt = dt.replace(year=datetime.now().year)
            timestamp_utc = dt.strftime("%H:%M UTC")
        except:
            timestamp_utc = timestamp

        # Mask part of the IP address
        ip_parts = ip.split('.')
        masked_ip = f"{ip_parts[0]}.{ip_parts[1]}.xx.xx"

        self.attempts.append({
            'timestamp': timestamp_utc,
            'username': username,
            'ip': masked_ip,
            'status': status,
            'raw_line': line.strip()  # Store raw line for debugging
        })

    def _extract_timestamp(self, line: str) -> str:
        """Extract timestamp from log line."""
        timestamp_match = re.search(r"(\w{3}\s+\d{1,2}\s+\d{2}:\d{2}:\d{2})", line)
        return timestamp_match.group(1) if timestamp_match else "Unknown"

    def _extract_ip(self, line: str) -> str:
        """Extract IP address from log line."""
        ip_match = re.search(r"from (\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})", line)
        return ip_match.group(1) if ip_match else None

    def _extract_username(self, line: str) -> str:
        """Extract username from log line."""
        if "invalid user" in line.lower():
            username_match = re.search(r"invalid user (\S+)", line, re.IGNORECASE)
        else:
            username_match = re.search(r"for (\S+) from", line)
        return username_match.group(1) if username_match else None

    def generate_report(self) -> str:
        """Generate a report grouped by IP address."""
        report = []
        report.append("=== SSH Connection Attempts Analysis Report ===")
        report.append(f"Log File: {self.log_path}")
        report.append(f"Analysis Date: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n")

        # Group attempts by IP
        ip_groups = defaultdict(list)
        for attempt in self.attempts:
            ip_groups[attempt['ip']].append(attempt)

        # Sort IPs by their first attempt timestamp
        sorted_ips = sorted(ip_groups.items(), 
                          key=lambda x: min(a['timestamp'] for a in x[1]))

        # Process each IP group
        for ip, attempts in sorted_ips:
            report.append(f"IP Address: {ip}")
            report.append("-" * 50)

            # Group attempts by username for this IP
            user_attempts = defaultdict(list)

            for attempt in attempts:
                user_attempts[attempt['username']].append(attempt)
                self.user_attempts[attempt['username']].append(attempt)
            # Sort users by their first attempt timestamp
            sorted_users = sorted(user_attempts.items(), 
                                key=lambda x: min(a['timestamp'] for a in x[1]))

            # Process each user's attempts
            for username, user_data in sorted_users:
                # Sort attempts chronologically
                sorted_attempts = sorted(user_data, key=lambda x: x['timestamp'])
                
                report.append(f"\nTime: {sorted_attempts[0]['timestamp']}")
                
                if len(sorted_attempts) > 1:
                    report.append(f"Last Attempt: {sorted_attempts[-1]['timestamp']}")
                    report.append(f"Total Attempts: {len(sorted_attempts)}")
                
                report.append(f"User: {username}")
                
                # Group attempts by status
                status_counts = defaultdict(int)
                for attempt in sorted_attempts:
                    status_counts[attempt['status']] += 1
                
                # Report each status type
                for status, count in status_counts.items():
                    if count == 1:
                        report.append(f"Status: {status}")
                    else:
                        report.append(f"Status: {status} ({count} attempts)")

            report.append("\n...")

        return "\n".join(report)
    
    
    def getUsernameAttempts(self,username)-> int:
      attempts = self.user_attempts[username]
      # Total number of attempts
      return len(attempts)

    def generate_report_by_username(self, username) -> str:
      # Check if the username exists in user_attempts
      if username not in self.user_attempts:
          return f"No login attempts found for user: {username}"

      # Get the user's login attempts
      attempts = self.user_attempts[username]

      # Total number of attempts
      total_attempts = len(attempts)

      # Get the last attempt's timestamp
      last_attempt_time = attempts[-1]['timestamp']

      # Count occurrences of each status
      status_counts = {}
      for attempt in attempts:
          status = attempt['status']
          status_counts[status] = status_counts.get(status, 0) + 1

      # Build the status report
      status_report = ", ".join([f"{status} ({count} attempts)" for status, count in status_counts.items()])

      # Print the report
      report = (
          f"Time: {attempts[0]['timestamp']}\n"
          f"Last Attempt: {last_attempt_time}\n"
          f"Total Attempts: {total_attempts}\n"
          f"User: {username}\n"
          f"Status: {status_report}"
      )
      return report
def send_notification(title, message):
    """
    Sends a system notification using plyer.
    
    :param title: Title of the notification
    :param message: Body of the notification
    """
    notification.notify(
        title=title,
        message=message,
        app_name="Auth script",
        timeout=10  # Notification duration in seconds
    )

def notify_attempts(path,username, attempts):
    while True:
        analyzer=SSHLogAnalyzer(path)
        analyzer.parse_log_file()
        total_attempts=analyzer.getUsernameAttempts(username)
        if attempts>=total_attempts:
            send_notification("Attemtps alert!",f"User {username} has reached {attempts} attempts")
        time.sleep(20)
def main():
    parser = argparse.ArgumentParser(description='Analyze SSH log files for connection attempts.')
    parser.add_argument('-u', '--user', help='Get logs for specific user')
    
    parser.add_argument('log_file', help='Path to the SSH log file')
    parser.add_argument('-o', '--output', help='Output file for the report')
    parser.add_argument('-d', '--debug', action='store_true', 
                       help='Include raw log lines in output')
    parser.add_argument('--attempts', help='Output file for the report')
    
    args = parser.parse_args()

   
    report=""
    if args.user:
        if args.attempts:
          if not args.attempts.isdigit():
              print("You must enter digit for the number of attempts")
              exit(-1)
              
          notify_attempts(args.log_file,args.user.strip(),int(args.attempts))
        # print(analyzer.generate_report_by_username("admin"))
        # report=analyzer.generate_report_by_username(f"{args.user}".strip())
        analyzer = SSHLogAnalyzer(args.log_file)
        analyzer.parse_log_file()
        analyzer.generate_report()
        report=analyzer.generate_report_by_username(args.user.strip())
        print()
        # report=analyzer.generate_report_by_username("admin")
      
    else:
        analyzer = SSHLogAnalyzer(args.log_file)
        analyzer.parse_log_file()
        report = analyzer.generate_report()

    if args.output:
        with open(args.output, 'w') as f:
            f.write(report)
        print(f"Report has been written to {args.output}")
    else:
        print(report)

if __name__ == "__main__":
    main()